package rabbitmq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
	"time"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *rabbitMQConnection) GetConnection() (*amqp.Connection, error) {
	// use a mutex lock around getting the connection because we don't want to accidentally have leaking connections
	//	in case two functions try to instantiate new connections at the same time
	r.mutex.Lock()
	if r.conn != nil && !r.conn.IsClosed() {
		r.mutex.Unlock()
		return r.conn, nil
	} else {
		for {
			logging.LogInfo("Attempting to connect to rabbitmq")
			conn, err := amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
				utils.MythicConfig.RabbitmqUser,
				utils.MythicConfig.RabbitmqPassword,
				utils.MythicConfig.RabbitmqHost,
				utils.MythicConfig.RabbitmqPort,
				utils.MythicConfig.RabbitmqVHost),
				amqp.Config{
					Dial: func(network, addr string) (net.Conn, error) {
						return net.DialTimeout(network, addr, 10*time.Second)
					},
				},
			)
			if err != nil {
				logging.LogError(err, "Failed to connect to rabbitmq")
				time.Sleep(RETRY_CONNECT_DELAY)
				continue
			}
			r.conn = conn
			r.mutex.Unlock()
			return conn, nil
		}
	}
}
func (r *rabbitMQConnection) SendStructMessage(exchange string, queue string, correlationId string, body interface{}) error {
	if jsonBody, err := json.Marshal(body); err != nil {
		return err
	} else {
		return r.SendMessage(exchange, queue, correlationId, jsonBody)
	}
}
func (r *rabbitMQConnection) SendRPCStructMessage(exchange string, queue string, body interface{}) ([]byte, error) {
	if inputBytes, err := json.Marshal(body); err != nil {
		logging.LogError(err, "Failed to convert input to JSON", "input", body)
		return nil, err
	} else {
		return r.SendRPCMessage(exchange, queue, inputBytes, true)
	}
}
func (r *rabbitMQConnection) SendMessage(exchange string, queue string, correlationId string, body []byte) error {
	// to send a normal message out to a direct queue set:
	// exchange: MYTHIC_EXCHANGE
	// queue: which routing key is listening (this is the direct name)
	// correlation_id: empty string
	if conn, err := r.GetConnection(); err != nil {
		return err
	} else if ch, err := conn.Channel(); err != nil {
		logging.LogError(err, "Failed to open rabbitmq channel")
		return err
	} else if err := ch.Confirm(false); err != nil {
		logging.LogError(err, "Channel could not be put into confirm mode")
		return err
	} else {
		msg := amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			Body:          body,
			DeliveryMode:  amqp.Persistent,
		}
		if err = ch.Publish(
			exchange, // exchange
			queue,    // routing key
			true,     // mandatory
			false,    // immediate
			msg,      // publishing
		); err != nil {
			logging.LogError(err, "there was an error publishing a message", "queue", queue)
			return err
		}
		select {
		case ntf := <-ch.NotifyPublish(make(chan amqp.Confirmation, 1)):
			if !ntf.Ack {
				err := errors.New("Failed to deliver message, not ACK-ed by receiver")
				logging.LogError(err, "failed to deliver message to exchange/queue, notifyPublish")
				return err
			}
		case ret := <-ch.NotifyReturn(make(chan amqp.Return)):
			err := errors.New(getMeaningfulRabbitmqError(ret))
			logging.LogError(err, "failed to deliver message to exchange/queue, NotifyReturn", "errorCode", ret.ReplyCode, "errorText", ret.ReplyText)
			return err
		case <-time.After(RPC_TIMEOUT):
			err := errors.New("Message delivery confirmation timed out")
			logging.LogError(err, "message delivery confirmation to exchange/queue timed out")
			return err
		}
		return nil
	}

}
func (r *rabbitMQConnection) SendRPCMessage(exchange string, queue string, body []byte, exclusiveQueue bool) ([]byte, error) {
	if conn, err := r.GetConnection(); err != nil {
		return nil, err
	} else if ch, err := conn.Channel(); err != nil {
		logging.LogError(err, "Failed to open rabbitmq channel")
		return nil, err
	} else if err := ch.Confirm(false); err != nil {
		logging.LogError(err, "Channel could not be put into confirm mode")
		return nil, err
	} else if err = ch.ExchangeDeclare(
		exchange, // exchange name
		"direct", // type of exchange, ex: topic, fanout, direct, etc
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		logging.LogError(err, "Failed to declare exchange", "exchange", exchange, "exchange_type", "direct", "retry_wait_time", RETRY_CONNECT_DELAY)
		return nil, err
	} else if msgs, err := ch.Consume(
		"amq.rabbitmq.reply-to", // queue name
		"",                      // consumer
		true,                    // auto-ack
		exclusiveQueue,          // exclusive
		false,                   // no local
		false,                   // no wait
		nil,                     // args
	); err != nil {
		logging.LogError(err, "Failed to start consuming for RPC replies")
		return nil, err
	} else {
		msg := amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: uuid.NewString(),
			Body:          body,
			ReplyTo:       "amq.rabbitmq.reply-to",
		}
		if err = ch.Publish(
			exchange, // exchange
			queue,    // routing key
			true,     // mandatory
			false,    // immediate
			msg,      // publishing
		); err != nil {
			logging.LogError(err, "there was an error publishing a message", "queue", queue)
			return nil, err
		}
		select {
		case ntf := <-ch.NotifyPublish(make(chan amqp.Confirmation, 1)):
			if !ntf.Ack {
				err := errors.New("Failed to deliver message, not ACK-ed by receiver")
				logging.LogError(err, "failed to deliver message to exchange/queue, notifyPublish")
				return nil, err
			}
		case ret := <-ch.NotifyReturn(make(chan amqp.Return)):
			err := errors.New(getMeaningfulRabbitmqError(ret))
			logging.LogError(err, "failed to deliver message to exchange/queue, NotifyReturn", "errorCode", ret.ReplyCode, "errorText", ret.ReplyText)
			return nil, err
		case <-time.After(RPC_TIMEOUT):
			err := errors.New("Message delivery confirmation timed out")
			logging.LogError(err, "message delivery confirmation to exchange/queue timed out")
			return nil, err
		}
		logging.LogDebug("Sent RPC message", "queue", queue)
		select {
		case m := <-msgs:
			logging.LogDebug("Got RPC Reply", "queue", queue)
			return m.Body, nil
		case <-time.After(RPC_TIMEOUT):
			logging.LogError(nil, "Timeout reached waiting for RPC reply")
			return nil, errors.New("Timeout reached waiting for RPC reply")
		}
	}

	/*
		_, err = ch.QueueDeclarePassive(
			queue, // name, queue
			false, // durable
			true,  // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			logging.LogError(err, "Failed to declare queue, RPC endpoint doesn't exist", "retry_wait_time", RETRY_CONNECT_DELAY)
			return nil, err
		}*/
}
func (r *rabbitMQConnection) ReceiveFromMythicDirectExchange(exchange string, queue string, routingKey string, handler QueueHandler, exclusiveQueue bool) {
	// exchange is a direct exchange
	// queue is where the messages get sent to (local name)
	// routingKey is the specific direct topic we're interested in for the exchange
	// handler processes the messages we get on our queue
	for {
		if conn, err := r.GetConnection(); err != nil {
			logging.LogError(err, "Failed to connect to rabbitmq", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if ch, err := conn.Channel(); err != nil {
			logging.LogError(err, "Failed to open rabbitmq channel", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if err = ch.ExchangeDeclare(
			exchange, // exchange name
			"direct", // type of exchange, ex: topic, fanout, direct, etc
			true,     // durable
			true,     // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		); err != nil {
			logging.LogError(err, "Failed to declare exchange", "exchange", exchange, "exchange_type", "direct", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if q, err := ch.QueueDeclare(
			queue,          // name, queue
			false,          // durable
			true,           // delete when unused
			exclusiveQueue, // exclusive
			false,          // no-wait
			nil,            // arguments
		); err != nil {
			logging.LogError(err, "Failed to declare queue", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if err = ch.QueueBind(
			q.Name,     // queue name
			routingKey, // routing key
			exchange,   // exchange name
			false,      // nowait
			nil,        // arguments
		); err != nil {
			logging.LogError(err, "Failed to bind to queue to receive messages", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if msgs, err := ch.Consume(
			q.Name, // queue name
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no local
			false,  // no wait
			nil,    // args
		); err != nil {
			logging.LogError(err, "Failed to start consuming on queue", "queue", q.Name)
		} else {
			forever := make(chan bool)
			go func() {
				for d := range msgs {
					go handler(d.Body)
				}
				forever <- true
			}()
			logging.LogInfo("Started listening for messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
			<-forever
			logging.LogError(nil, "Stopped listening for messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
		}

	}
}
func (r *rabbitMQConnection) ReceiveFromMythicDirectTopicExchange(exchange string, queue string, routingKey string, handler QueueHandler, exclusiveQueue bool) {
	// exchange is a direct exchange
	// queue is where the messages get sent to (local name)
	// routingKey is the specific direct topic we're interested in for the exchange
	// handler processes the messages we get on our queue
	for {
		if conn, err := r.GetConnection(); err != nil {
			logging.LogError(err, "Failed to connect to rabbitmq", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if ch, err := conn.Channel(); err != nil {
			logging.LogError(err, "Failed to open rabbitmq channel", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if err = ch.ExchangeDeclare(
			exchange, // exchange name
			"topic",  // type of exchange, ex: topic, fanout, direct, etc
			true,     // durable
			true,     // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		); err != nil {
			logging.LogError(err, "Failed to declare exchange", "exchange", exchange, "exchange_type", "direct", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if q, err := ch.QueueDeclare(
			"",             // name, queue
			false,          // durable
			true,           // delete when unused
			exclusiveQueue, // exclusive
			false,          // no-wait
			nil,            // arguments
		); err != nil {
			logging.LogError(err, "Failed to declare queue", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if err = ch.QueueBind(
			q.Name,     // queue name
			routingKey, // routing key
			exchange,   // exchange name
			false,      // nowait
			nil,        // arguments
		); err != nil {
			logging.LogError(err, "Failed to bind to queue to receive messages", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if msgs, err := ch.Consume(
			q.Name, // queue name
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no local
			false,  // no wait
			nil,    // args
		); err != nil {
			logging.LogError(err, "Failed to start consuming on queue", "queue", q.Name)
		} else {
			forever := make(chan bool)
			go func() {
				for d := range msgs {
					go handler(d.Body)
				}
				forever <- true
			}()
			logging.LogInfo("Started listening for messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
			<-forever
			logging.LogError(nil, "Stopped listening for messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
		}

	}
}
func (r *rabbitMQConnection) ReceiveFromRPCQueue(exchange string, queue string, routingKey string, handler RPCQueueHandler, exclusiveQueue bool) {
	for {
		if conn, err := r.GetConnection(); err != nil {
			logging.LogError(err, "Failed to connect to rabbitmq", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if ch, err := conn.Channel(); err != nil {
			logging.LogError(err, "Failed to open rabbitmq channel", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if err = ch.ExchangeDeclare(
			exchange, // exchange name
			"direct", // type of exchange, ex: topic, fanout, direct, etc
			true,     // durable
			true,     // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		); err != nil {
			logging.LogError(err, "Failed to declare exchange", "exchange", exchange, "exchange_type", "direct", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if q, err := ch.QueueDeclare(
			queue,          // name, queue
			false,          // durable
			true,           // delete when unused
			exclusiveQueue, // exclusive
			false,          // no-wait
			nil,            // arguments
		); err != nil {
			logging.LogError(err, "Failed to declare queue", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if err = ch.QueueBind(
			q.Name,     // queue name
			routingKey, // routing key
			exchange,   // exchange name
			false,      // nowait
			nil,        // arguments
		); err != nil {
			logging.LogError(err, "Failed to bind to queue to receive messages", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		} else if msgs, err := ch.Consume(
			q.Name,         // queue name
			"",             // consumer
			false,          // auto-ack
			exclusiveQueue, // exclusive
			false,          // no local
			false,          // no wait
			nil,            // args
		); err != nil {
			logging.LogError(err, "Failed to start consuming messages on queue", "queue", q.Name)
			continue
		} else {
			forever := make(chan bool)
			go func() {
				for d := range msgs {
					responseMsg := handler(d.Body)
					if responseMsgJson, err := json.Marshal(responseMsg); err != nil {
						logging.LogError(err, "Failed to generate JSON for getFile response")
						continue
					} else if err = ch.Publish(
						"",        // exchange
						d.ReplyTo, //routing key
						true,      // mandatory
						false,     // immediate
						amqp.Publishing{
							ContentType:   "application/json",
							Body:          responseMsgJson,
							CorrelationId: d.CorrelationId,
							DeliveryMode:  amqp.Persistent,
						}); err != nil {
						logging.LogError(err, "Failed to send message")
					} else if err = ch.Ack(d.DeliveryTag, false); err != nil {
						logging.LogError(err, "Failed to Ack message")
					}
				}
				forever <- true
			}()
			logging.LogInfo("Started listening for rpc messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
			<-forever
			logging.LogError(nil, "Stopped listening for messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
		}

	}
}
func (r *rabbitMQConnection) CheckConsumerExists(exchange string, queue string, exclusiveQueue bool) (bool, error) {
	//logging.LogDebug("checking queue existence", "queue", queue)
	conn, err := r.GetConnection()
	if err != nil {
		return false, err
	}
	ch, err := conn.Channel()
	if err != nil {
		logging.LogError(err, "Failed to open rabbitmq channel")
		return false, err
	}
	defer ch.Close()
	if err := ch.Confirm(false); err != nil {
		logging.LogError(err, "Channel could not be put into confirm mode")
		return false, err
	}
	err = ch.ExchangeDeclare(
		exchange, // exchange name
		"direct", // type of exchange, ex: topic, fanout, direct, etc
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logging.LogError(err, "Failed to declare exchange", "exchange", exchange, "exchange_type", "direct", "retry_wait_time", RETRY_CONNECT_DELAY)
		return false, err
	}

	if _, err := ch.QueueDeclarePassive(
		queue,          // name, queue
		false,          // durable
		true,           // delete when unused
		exclusiveQueue, // exclusive
		false,          // no-wait
		nil,            // arguments
	); err != nil {
		errorMessage := err.Error()
		//logging.LogError(err, "Error when checking for queue")
		if strings.Contains(errorMessage, "Exception (405)") {
			return true, nil
		} else if strings.Contains(errorMessage, "Exception (404)") {
			return false, nil
		} else {
			logging.LogError(err, "Unknown error (not 404 or 405) when checking for container existence")
			return false, err
		}
	} else {
		return true, nil
	}
}
func getMeaningfulRabbitmqError(ret amqp.Return) string {
	switch ret.ReplyCode {
	case 312:
		return fmt.Sprintf("No RabbitMQ Route for %s. Is the container online (./mythic-cli status)?\nIf the container is online, there might be an issue within the container processing the request (./mythic-cli logs [container name]). ", ret.RoutingKey)
	default:
		return fmt.Sprintf("Failed to deliver message to exchange/queue. Error code: %d, Error Text: %s", ret.ReplyCode, ret.ReplyText)
	}
}

// payload helper functions
func UploadPayloadData(payloadBuildMsg agentstructs.PayloadBuildMessage, payloadBuildMsgResponse agentstructs.PayloadBuildResponse) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if fileWriter, err := writer.CreateFormFile("file", "payload"); err != nil {
		logging.LogError(err, "Failed to create new form file to upload payload")
		return err
	} else if _, err = io.Copy(fileWriter, bytes.NewReader(*payloadBuildMsgResponse.Payload)); err != nil {
		logging.LogError(err, "Failed to write payload bytes to form")
		return err
	} else if fieldWriter, err := writer.CreateFormField("agent-file-id"); err != nil {
		logging.LogError(err, "Failed to add new form field to upload payload")
		return err
	} else if _, err := fieldWriter.Write([]byte(payloadBuildMsg.PayloadFileUUID)); err != nil {
		logging.LogError(err, "Failed to add in agent-file-id to form")
		return err
	}
	writer.Close()
	if request, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/direct/upload/%s", utils.MythicConfig.MythicServerHost, utils.MythicConfig.MythicServerPort, payloadBuildMsg.PayloadFileUUID), body); err != nil {
		logging.LogError(err, "Failed to create new POST request to send payload to Mythic")
		return err
	} else {
		request.Header.Add("Content-Type", writer.FormDataContentType())
		if resp, err := http.DefaultClient.Do(request); err != nil {
			logging.LogError(err, "Failed to send payload over to Mythic")
			return err
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				logging.LogError(nil, "Failed to send payload to Mythic", "status code", resp.StatusCode)
				return errors.New(fmt.Sprintf("Failed to send payload to Mythic with status code: %d\n", resp.StatusCode))
			}
		}
	}

	return nil
}
func WrapPayloadBuild(msg []byte) {
	payloadBuildMsg := agentstructs.PayloadBuildMessage{}
	if err := json.Unmarshal(msg, &payloadBuildMsg); err != nil {
		logging.LogError(err, "Failed to process payload build message")
	} else {
		var payloadBuildResponse agentstructs.PayloadBuildResponse
		if payloadBuildFunc := agentstructs.AllPayloadData.Get(payloadBuildMsg.PayloadType).GetBuildFunction(); payloadBuildFunc == nil {
			logging.LogError(nil, "Failed to get payload build function. Do you have a function called 'build'?")
			payloadBuildResponse.Success = false
		} else {
			payloadBuildResponse = agentstructs.AllPayloadData.Get(payloadBuildMsg.PayloadType).GetBuildFunction()(payloadBuildMsg)
		}
		// handle sending off the payload via a web request separately from the rest of the message
		if payloadBuildResponse.Payload != nil {
			if err := UploadPayloadData(payloadBuildMsg, payloadBuildResponse); err != nil {
				logging.LogError(err, "Failed to send payload back to Mythic via web request")
				payloadBuildResponse.BuildMessage = payloadBuildResponse.BuildMessage + "\nFailed to send payload back to Mythic: " + err.Error()
				payloadBuildResponse.Success = false
			}
		}
		if err := RabbitMQConnection.SendStructMessage(
			MYTHIC_EXCHANGE,
			PT_BUILD_RESPONSE_ROUTING_KEY,
			"",
			payloadBuildResponse,
		); err != nil {
			logging.LogError(err, "Failed to send payload response back to Mythic")
		}
		logging.LogDebug("Finished processing payload build message")
	}
}
func prepTaskArgs(command agentstructs.Command, taskMessage *agentstructs.PTTaskMessageAllData) error {
	if args, err := agentstructs.GenerateArgsData(command.CommandParameters, *taskMessage); err != nil {
		logging.LogError(err, "Failed to generate args data for create tasking")
		return errors.New(fmt.Sprintf("Failed to generate args data:\n%s", err.Error()))
	} else {
		taskMessage.Args = args
		switch taskMessage.Args.GetTaskingLocation() {
		case "command_line":
			if command.TaskFunctionParseArgString != nil {
				// from scripting or if there are no parameters defined, this gets called
				if err := command.TaskFunctionParseArgString(&taskMessage.Args, taskMessage.Args.GetCommandLine()); err != nil {
					logging.LogError(err, "Failed to run ParseArgString function", "command Name", command.Name)
					return errors.New(fmt.Sprintf("Failed to run %s's ParseArgString function:\n%s", command.Name, err.Error()))
				}
			}
		default:
			// try to parse dictionary first - parsed_cli, browser script, etc
			if command.TaskFunctionParseArgDictionary != nil {
				tempArgs := map[string]interface{}{}
				if err := json.Unmarshal([]byte(taskMessage.Args.GetCommandLine()), &tempArgs); err != nil {
					// failed to parse as a dictionary, so try parsing as
					if command.TaskFunctionParseArgString != nil {
						if err := command.TaskFunctionParseArgString(&taskMessage.Args, taskMessage.Args.GetCommandLine()); err != nil {
							logging.LogError(err, "Failed to run ParseArgString function", "command Name", command.Name)
							return errors.New(fmt.Sprintf("Failed to run %s's ParseArgString function:\n%s", command.Name, err.Error()))
						} else {
							break
						}
					}
					logging.LogError(err, "Failed to parse arguments from parsed_cli into dictionary")
					return errors.New(fmt.Sprintf("Failed to parse arguments from parsed_cli into dictionary:\n%s", err.Error()))
				} else if err := command.TaskFunctionParseArgDictionary(&taskMessage.Args, tempArgs); err != nil {
					logging.LogError(err, "Failed to run ParseArgDictionary function", "command Name", command.Name)
					return errors.New(fmt.Sprintf("Failed to run %s's ParseArgDictionary function:\n%s", command.Name, err.Error()))
				}
				// if no dictionary function, fall back to the string function if it exists
			} else if command.TaskFunctionParseArgString != nil {

				if err := command.TaskFunctionParseArgString(&taskMessage.Args, taskMessage.Args.GetCommandLine()); err != nil {
					logging.LogError(err, "Failed to run ParseArgString function", "command Name", command.Name)
					return errors.New(fmt.Sprintf("Failed to run %s's ParseArgString function:\n%s", command.Name, err.Error()))
				}
			}
		}
		// before we process the function's create_tasking, let's make sure that we have all the right arguments supplied by the user
		// if something is "required" then the user needs to specify it, otherwise it's not required because a default value works fine
		if requiredArgsHaveValues, err := taskMessage.Args.VerifyRequiredArgsHaveValues(); err != nil {
			logging.LogError(err, "Failed to verify if all required args have values")
			return errors.New(fmt.Sprintf("Failed to verify if all required args have values:\n%s", err.Error()))
		} else if !requiredArgsHaveValues {
			return errors.New(fmt.Sprintf("Some required args are missing values"))
		}
		return nil
	}
}

package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	c2structs "github.com/MythicMeta/MythicContainer/c2_structs"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/config"
	"github.com/MythicMeta/MythicContainer/logging"
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
		if config.MythicConfig.RabbitmqHost == "" {
			log.Fatalf("[-] Missing RABBITMQ_HOST environment variable point to rabbitmq server IP")
		}
		for {
			logging.LogInfo("Attempting to connect to rabbitmq")
			conn, err := amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
				config.MythicConfig.RabbitmqUser,
				config.MythicConfig.RabbitmqPassword,
				config.MythicConfig.RabbitmqHost,
				config.MythicConfig.RabbitmqPort,
				config.MythicConfig.RabbitmqVHost),
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
func (r *rabbitMQConnection) SendStructMessage(exchange string, queue string, correlationId string, body interface{}, ignoreErrorMessage bool) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return r.SendMessage(exchange, queue, correlationId, jsonBody, ignoreErrorMessage)
}
func (r *rabbitMQConnection) SendRPCStructMessage(exchange string, queue string, body interface{}) ([]byte, error) {
	inputBytes, err := json.Marshal(body)
	if err != nil {
		logging.LogError(err, "Failed to convert input to JSON", "input", body)
		return nil, err
	}
	return r.SendRPCMessage(exchange, queue, inputBytes, true)
}
func (r *rabbitMQConnection) SendMessage(exchange string, queue string, correlationId string, body []byte, ignoreErrormessage bool) error {
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
		ch.Close()
		return err
	} else {
		defer ch.Close()
		msg := amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			Body:          body,
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
			if !ignoreErrormessage {
				logging.LogError(err, "failed to deliver message to exchange/queue, NotifyReturn", "errorCode", ret.ReplyCode, "errorText", ret.ReplyText)
			}
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
	conn, err := r.GetConnection()
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		logging.LogError(err, "Failed to open rabbitmq channel")
		return nil, err
	}
	err = ch.Confirm(false)
	if err != nil {
		logging.LogError(err, "Channel could not be put into confirm mode")
		ch.Close()
		return nil, err
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
		return nil, err
	}
	msgs, err := ch.Consume(
		"amq.rabbitmq.reply-to", // queue name
		"",                      // consumer
		true,                    // auto-ack
		exclusiveQueue,          // exclusive
		false,                   // no local
		false,                   // no wait
		nil,                     // args
	)
	if err != nil {
		logging.LogError(err, "Failed to start consuming for RPC replies")
		ch.Close()
		return nil, err
	}
	defer ch.Close()
	msg := amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: uuid.NewString(),
		Body:          body,
		ReplyTo:       "amq.rabbitmq.reply-to",
	}
	for attempt := 0; attempt < 3; attempt++ {
		err = ch.Publish(
			exchange, // exchange
			queue,    // routing key
			true,     // mandatory
			false,    // immediate
			msg,      // publishing
		)
		if err != nil {
			logging.LogError(err, "there was an error publishing a message, trying again", "queue", queue)
			continue
		}
		select {
		case ntf := <-ch.NotifyPublish(make(chan amqp.Confirmation, 1)):
			if !ntf.Ack {
				err = errors.New("Failed to deliver message, not ACK-ed by receiver")
				logging.LogError(err, "failed to deliver message to exchange/queue, notifyPublish, trying again", "queue", queue)
				continue
			}
		case ret := <-ch.NotifyReturn(make(chan amqp.Return, 1)):
			err = errors.New(getMeaningfulRabbitmqError(ret))
			continue
		case <-time.After(RPC_TIMEOUT):
			err = errors.New("message delivery confirmation timed out in SendRPCMessage")
			logging.LogError(err, "message delivery confirmation to exchange/queue timed out, trying again", "queue", queue)
			continue
		}
		//logging.LogDebug("Sent RPC message", "queue", queue)
		select {
		case m := <-msgs:
			//logging.LogDebug("Got RPC Reply", "queue", queue)
			return m.Body, nil
		case <-time.After(RPC_TIMEOUT):
			logging.LogError(nil, "Timeout reached waiting for RPC reply, trying again", "queue", queue)
			err = errors.New("Timeout reached waiting for RPC reply")
			continue
		}
	}
	logging.LogError(err, "failed 3 times")
	return nil, err
}
func (r *rabbitMQConnection) ReceiveFromMythicDirectExchange(exchange string, queue string, routingKey string, handler QueueHandler, exclusiveQueue bool, wg *sync.WaitGroup) {
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
			ch.Close()
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
			ch.Close()
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
			ch.Close()
		} else {
			forever := make(chan bool)
			go func() {
				for d := range msgs {
					go handler(d.Body)
				}
				forever <- true
			}()
			logging.LogInfo("Started listening for messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
			if wg != nil {
				wg.Done()
				wg = nil
			}
			<-forever
			ch.Close()
			logging.LogError(nil, "Stopped listening for messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
		}

	}
}
func (r *rabbitMQConnection) ReceiveFromRPCQueue(exchange string, queue string, routingKey string, handler RPCQueueHandler, exclusiveQueue bool, wg *sync.WaitGroup) {
	for {
		conn, err := r.GetConnection()
		if err != nil {
			logging.LogError(err, "Failed to connect to rabbitmq", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		}
		ch, err := conn.Channel()
		if err != nil {
			logging.LogError(err, "Failed to open rabbitmq channel", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
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
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		}
		q, err := ch.QueueDeclare(
			queue,          // name, queue
			true,           // durable
			false,          // delete when unused
			exclusiveQueue, // exclusive
			false,          // no-wait
			nil,            // arguments
		)
		if err != nil {
			logging.LogError(err, "Failed to declare queue", "retry_wait_time", RETRY_CONNECT_DELAY)
			ch.Close()
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		}
		err = ch.QueueBind(
			q.Name,     // queue name
			routingKey, // routing key
			exchange,   // exchange name
			false,      // nowait
			nil,        // arguments
		)
		if err != nil {
			logging.LogError(err, "Failed to bind to queue to receive messages", "retry_wait_time", RETRY_CONNECT_DELAY)
			time.Sleep(RETRY_CONNECT_DELAY)
			ch.Close()
			continue
		}
		msgs, err := ch.Consume(
			q.Name,         // queue name
			"",             // consumer
			false,          // auto-ack
			exclusiveQueue, // exclusive
			false,          // no local
			false,          // no wait
			nil,            // args
		)
		if err != nil {
			logging.LogError(err, "Failed to start consuming messages on queue", "queue", q.Name)
			ch.Close()
			continue
		}
		forever := make(chan bool)
		go func() {
			for d := range msgs {
				//logging.LogInfo("about to handle rpc msg", "queue", q.Name)
				responseMsg := handler(d.Body)
				//logging.LogInfo("finished handling rpc msg", "queue", q.Name)
				responseMsgJson, err := json.Marshal(responseMsg)
				if err != nil {
					logging.LogError(err, "Failed to generate JSON for getFile response")
					continue
				}
				err = ch.Publish(
					"",        // exchange
					d.ReplyTo, //routing key
					true,      // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "application/json",
						Body:          responseMsgJson,
						CorrelationId: d.CorrelationId,
					})
				if err != nil {
					logging.LogError(err, "Failed to send message")
					continue
				}
				err = ch.Ack(d.DeliveryTag, false)
				if err != nil {
					logging.LogError(err, "Failed to Ack message")
				}
			}
			forever <- true
		}()
		logging.LogInfo("Started listening for rpc messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
		if wg != nil {
			wg.Done()
			wg = nil
		}
		<-forever
		ch.Close()
		logging.LogError(nil, "Stopped listening for messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
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
func (r *rabbitMQConnection) GetNumberOfConsumersDirectChannels(exchange string, kind string, queue string) (uint, error) {
	//logging.LogDebug("checking queue existence", "queue", queue)
	if conn, err := r.GetConnection(); err != nil {
		return 0, err
	} else if ch, err := conn.Channel(); err != nil {
		logging.LogError(err, "Failed to open rabbitmq channel")
		return 0, err
	} else if err = ch.Confirm(false); err != nil {
		logging.LogError(err, "Channel could not be put into confirm mode")
		return 0, err
	} else if err = ch.ExchangeDeclare(
		exchange, // exchange name
		kind,     // type of exchange, ex: topic, fanout, direct, etc
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		logging.LogError(err, "Failed to declare exchange", "exchange", MYTHIC_TOPIC_EXCHANGE, "exchange_type", "topic", "retry_wait_time", RETRY_CONNECT_DELAY)
		return 0, err

	} else if q, err := ch.QueueDeclare(
		queue, // name, queue
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		logging.LogError(err, "Unknown error (not 404 or 405) when checking for container existence")
		return 0, err
	} else if err = ch.QueueBind(
		q.Name,   // queue name
		queue,    // routing key
		exchange, // exchange name
		false,    // nowait
		nil,      // arguments
	); err != nil {
		logging.LogError(err, "Failed to bind to queue to receive messages", "retry_wait_time", RETRY_CONNECT_DELAY)
		return 0, err
	} else {
		return uint(q.Consumers), nil
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
			"",             // name, queue - this needs to be unique or we round robin amongst listeners
			false,          // durable
			true,           // delete when unused
			exclusiveQueue, // exclusive
			false,          // no-wait
			nil,            // arguments
		); err != nil {
			logging.LogError(err, "Failed to declare queue", "retry_wait_time", RETRY_CONNECT_DELAY)
			ch.Close()
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
			ch.Close()
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
			ch.Close()
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
			ch.Close()
			logging.LogError(nil, "Stopped listening for messages", "exchange", exchange, "queue", queue, "routingKey", routingKey)
		}

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

func prepTaskArgs(command agentstructs.Command, taskMessage *agentstructs.PTTaskMessageAllData) error {
	args, err := agentstructs.GenerateArgsData(command.CommandParameters, *taskMessage)
	if err != nil {
		logging.LogError(err, "Failed to generate args data for create tasking")
		return errors.New(fmt.Sprintf("Failed to generate args data:\n%s", err.Error()))
	}
	taskMessage.Args = args
	switch taskMessage.Args.GetTaskingLocation() {
	case "command_line":
		if command.TaskFunctionParseArgString != nil {
			// from scripting or if there are no parameters defined, this gets called
			err = command.TaskFunctionParseArgString(&taskMessage.Args, taskMessage.Args.GetCommandLine())
			if err != nil {
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
				} else {
					logging.LogError(err, "Failed to parse arguments from parsed_cli into dictionary and no ParseArgString function defined, using raw command line")
				}

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
	// in case we auto-parsed a typed array into ["", "val"] data, submit it for processing first
	for _, arg := range taskMessage.Args.GetTypedArrayEntriesThatNeedProcessing() {
		//logging.LogInfo("arg would need more processing", "arg", arg)
		currentArg, _ := taskMessage.Args.GetTypedArrayArg(arg.Name)
		newUntypedArray := make([]string, len(currentArg))
		for i := 0; i < len(newUntypedArray); i++ {
			newUntypedArray[i] = currentArg[i][1]
		}
		newArg := arg.TypedArrayParseFunction(agentstructs.PTRPCTypedArrayParseFunctionMessage{
			Command:       command.Name,
			ParameterName: arg.Name,
			PayloadType:   taskMessage.PayloadType,
			Callback:      taskMessage.Callback.ID,
			InputArray:    newUntypedArray,
		})
		err = taskMessage.Args.SetArgValue(arg.Name, newArg)
		if err != nil {
			logging.LogError(err, "failed to set new typed array value")
		}
	}
	// before we process the function's create_tasking, let's make sure that we have all the right arguments supplied by the user
	// if something is "required" then the user needs to specify it, otherwise it's not required because a default value works fine
	requiredArgsHaveValues, err := taskMessage.Args.VerifyRequiredArgsHaveValues()
	if err != nil {
		logging.LogError(err, "Failed to verify if all required args have values")
		return errors.New(fmt.Sprintf("Failed to verify if all required args have values:\n%s", err.Error()))
	}
	if !requiredArgsHaveValues {
		return errors.New(fmt.Sprintf("Some required args are missing values"))
	}
	return nil

}

func restartC2Server(name string) {
	stopResponse := C2RPCStopServer(c2structs.C2RPCStopServerMessage{
		Name: name,
	})
	if !stopResponse.Success {
		_, _ = SendMythicRPCC2UpdateStatus(MythicRPCC2UpdateStatusMessage{
			Error:                 stopResponse.Error,
			InternalServerRunning: stopResponse.InternalServerRunning,
			C2Profile:             name,
		})
		return
	}
	startResponse := C2RPCStartServer(c2structs.C2RPCStartServerMessage{
		Name: name,
	})
	_, _ = SendMythicRPCC2UpdateStatus(MythicRPCC2UpdateStatusMessage{
		Error:                 startResponse.Error,
		InternalServerRunning: startResponse.InternalServerRunning,
		C2Profile:             name,
	})
}

package rabbitmq

import (
	"context"
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

type rpcResponse struct {
	Body []byte
	Err  error
}

func (r *rabbitMQConnection) GetConnection() (*amqp.Connection, error) {
	// use a mutex lock around getting the connection because we don't want to accidentally have leaking connections
	//	in case two functions try to instantiate new connections at the same time
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.conn != nil && !r.conn.IsClosed() {
		return r.conn, nil
	}
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
		return conn, nil
	}
}
func (r *rabbitMQConnection) SendStructMessage(exchange string, queue string, correlationId string, body interface{}, ignoreErrorMessage bool) error {
	return r.SendStructMessageWithContext(context.Background(), exchange, queue, correlationId, body, ignoreErrorMessage)
}

func (r *rabbitMQConnection) SendStructMessageWithContext(ctx context.Context, exchange string, queue string, correlationId string, body interface{}, ignoreErrorMessage bool) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return r.SendMessageWithContext(ctx, exchange, queue, correlationId, jsonBody, ignoreErrorMessage)
}
func (r *rabbitMQConnection) SendRPCStructMessage(exchange string, queue string, body interface{}, retryPolicy ...RPCRetryPolicy) ([]byte, error) {
	return r.SendRPCStructMessageWithContext(context.Background(), exchange, queue, body, retryPolicy...)
}

func (r *rabbitMQConnection) SendRPCStructMessageWithContext(ctx context.Context, exchange string, queue string, body interface{}, retryPolicy ...RPCRetryPolicy) ([]byte, error) {
	inputBytes, err := json.Marshal(body)
	if err != nil {
		logging.LogError(err, "Failed to convert input to JSON", "input", body)
		return nil, err
	}
	return r.SendRPCMessageWithContext(ctx, exchange, queue, inputBytes, true, retryPolicy...)
}

func (r *rabbitMQConnection) getPublisherChannel() (*amqp.Channel, chan amqp.Confirmation, chan amqp.Return, error) {
	if r.publisherChannel != nil && !r.publisherChannel.IsClosed() {
		return r.publisherChannel, r.publisherConfirm, r.publisherReturn, nil
	}
	conn, err := r.GetConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}
	if err = ch.Confirm(false); err != nil {
		ch.Close()
		return nil, nil, nil, err
	}
	r.publisherChannel = ch
	r.publisherConfirm = ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	r.publisherReturn = ch.NotifyReturn(make(chan amqp.Return, 1))
	return r.publisherChannel, r.publisherConfirm, r.publisherReturn, nil
}

func (r *rabbitMQConnection) resetPublisherChannel(ch *amqp.Channel) {
	if ch != nil && !ch.IsClosed() {
		ch.Close()
	}
	if r.publisherChannel == ch {
		r.publisherChannel = nil
		r.publisherConfirm = nil
		r.publisherReturn = nil
	}
}

func (r *rabbitMQConnection) SendMessage(exchange string, queue string, correlationId string, body []byte, ignoreErrormessage bool) error {
	return r.SendMessageWithContext(context.Background(), exchange, queue, correlationId, body, ignoreErrormessage)
}

func (r *rabbitMQConnection) SendMessageWithContext(ctx context.Context, exchange string, queue string, correlationId string, body []byte, ignoreErrormessage bool) error {
	// to send a normal message out to a direct queue set:
	// exchange: MYTHIC_EXCHANGE
	// queue: which routing key is listening (this is the direct name)
	// correlation_id: empty string
	for attempt := 0; attempt < 3; attempt++ {
		r.publisherMutex.Lock()
		ch, confirmChannel, notifyReturnChannel, err := r.getPublisherChannel()
		if err != nil {
			logging.LogError(err, "Failed to get rabbitmq publisher channel", "queue", queue)
			r.publisherMutex.Unlock()
			time.Sleep(RPC_TIMEOUT)
			continue
		}
		msg := amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			Body:          body,
			Headers:       HeadersFromContext(ctx),
		}
		err = ch.Publish(
			exchange, // exchange
			queue,    // routing key
			true,     // mandatory
			false,    // immediate
			msg,      // publishing
		)
		if err != nil {
			logging.LogError(err, "there was an error publishing a message", "queue", queue)
			r.resetPublisherChannel(ch)
			r.publisherMutex.Unlock()
			time.Sleep(RPC_TIMEOUT)
			continue
		}
		select {
		case ntf := <-confirmChannel:
			if !ntf.Ack {
				err = errors.New("Failed to deliver message, not ACK-ed by receiver")
				logging.LogError(err, "failed to deliver message to exchange/queue, notifyPublish")
				r.resetPublisherChannel(ch)
				r.publisherMutex.Unlock()
				time.Sleep(RPC_TIMEOUT)
				continue
			}
		case ret := <-notifyReturnChannel:
			err = errors.New(getMeaningfulRabbitmqError(ret))
			if !ignoreErrormessage {
				logging.LogError(err, "failed to deliver message to exchange/queue, NotifyReturn", "errorCode", ret.ReplyCode, "errorText", ret.ReplyText)
			}
			r.resetPublisherChannel(ch)
			r.publisherMutex.Unlock()
			time.Sleep(RPC_TIMEOUT)
			continue
		case <-time.After(RPC_TIMEOUT):
			err = errors.New("Message delivery confirmation timed out")
			logging.LogError(err, "message delivery confirmation to exchange/queue timed out")
			r.resetPublisherChannel(ch)
			r.publisherMutex.Unlock()
			continue
		}
		r.publisherMutex.Unlock()
		return nil
	}
	if !ignoreErrormessage {
		logging.LogError(errors.New("failed 3 times"), "failed 3 times", "queue", queue)
	}
	return errors.New(fmt.Sprintf("failed 3 times to send to queue %s", queue))
}

func getRPCRetryPolicy(retryPolicy []RPCRetryPolicy) RPCRetryPolicy {
	if len(retryPolicy) > 0 {
		return retryPolicy[0]
	}
	return RPC_RETRY_POLICY_RETRY_ON_TIMEOUT
}

func (r *rabbitMQConnection) getRPCTimeout(retryPolicy RPCRetryPolicy) time.Duration {
	if retryPolicy == RPC_RETRY_POLICY_CUSTOM_TIMEOUT && config.MythicConfig.CustomRPCTimeout > 0 {
		return config.MythicConfig.CustomRPCTimeout
	}
	return RPC_TIMEOUT
}

func (r *rabbitMQConnection) getRPCClientLocked(exchange string, exclusiveQueue bool) (*amqp.Channel, chan amqp.Confirmation, chan amqp.Return, error) {
	if r.rpcChannel != nil && !r.rpcChannel.IsClosed() {
		if err := r.declareRPCExchangeLocked(r.rpcChannel, exchange); err != nil {
			r.resetRPCClientLocked(r.rpcChannel, err)
			return nil, nil, nil, err
		}
		return r.rpcChannel, r.rpcConfirm, r.rpcReturn, nil
	}
	conn, err := r.GetConnection()
	if err != nil {
		return nil, nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}
	if err = ch.Confirm(false); err != nil {
		ch.Close()
		return nil, nil, nil, err
	}
	confirmChannel := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	notifyReturnChannel := ch.NotifyReturn(make(chan amqp.Return, 1))
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
		ch.Close()
		return nil, nil, nil, err
	}
	r.rpcChannel = ch
	r.rpcConfirm = confirmChannel
	r.rpcReturn = notifyReturnChannel
	if r.rpcPending == nil {
		r.rpcPending = make(map[string]chan rpcResponse)
	}
	r.rpcExchanges = make(map[string]bool)
	if err = r.declareRPCExchangeLocked(ch, exchange); err != nil {
		r.resetRPCClientLocked(ch, err)
		return nil, nil, nil, err
	}
	go r.listenForRPCReplies(ch, msgs)
	return r.rpcChannel, r.rpcConfirm, r.rpcReturn, nil
}

func (r *rabbitMQConnection) declareRPCExchangeLocked(ch *amqp.Channel, exchange string) error {
	if r.rpcExchanges == nil {
		r.rpcExchanges = make(map[string]bool)
	}
	if r.rpcExchanges[exchange] {
		return nil
	}
	err := ch.ExchangeDeclare(
		exchange, // exchange name
		"direct", // type of exchange, ex: topic, fanout, direct, etc
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err == nil {
		r.rpcExchanges[exchange] = true
	}
	return err
}

func (r *rabbitMQConnection) listenForRPCReplies(ch *amqp.Channel, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		r.rpcClientMutex.Lock()
		responseChannel := r.rpcPending[d.CorrelationId]
		if responseChannel != nil {
			delete(r.rpcPending, d.CorrelationId)
		}
		r.rpcClientMutex.Unlock()
		if responseChannel != nil {
			responseChannel <- rpcResponse{Body: d.Body}
		}
	}
	r.rpcClientMutex.Lock()
	if r.rpcChannel == ch {
		r.resetRPCClientLocked(ch, errors.New("rpc reply consumer stopped"))
	}
	r.rpcClientMutex.Unlock()
}

func (r *rabbitMQConnection) resetRPCClientLocked(ch *amqp.Channel, err error) {
	if ch != nil && !ch.IsClosed() {
		ch.Close()
	}
	if r.rpcChannel == ch {
		for correlationID, responseChannel := range r.rpcPending {
			select {
			case responseChannel <- rpcResponse{Err: err}:
			default:
			}
			delete(r.rpcPending, correlationID)
		}
		r.rpcChannel = nil
		r.rpcConfirm = nil
		r.rpcReturn = nil
		r.rpcExchanges = nil
	}
}

func (r *rabbitMQConnection) removePendingRPCResponse(correlationID string) {
	r.rpcClientMutex.Lock()
	delete(r.rpcPending, correlationID)
	r.rpcClientMutex.Unlock()
}

func (r *rabbitMQConnection) publishRPCMessage(ctx context.Context, exchange string, queue string, correlationID string, body []byte, exclusiveQueue bool, responseChannel chan rpcResponse) error {
	r.rpcClientMutex.Lock()
	defer r.rpcClientMutex.Unlock()
	ch, confirmChannel, notifyReturnChannel, err := r.getRPCClientLocked(exchange, exclusiveQueue)
	if err != nil {
		logging.LogError(err, "Failed to get rabbitmq rpc channel", "queue", queue)
		return err
	}
	r.rpcPending[correlationID] = responseChannel
	msg := amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: correlationID,
		Body:          body,
		ReplyTo:       "amq.rabbitmq.reply-to",
		Headers:       HeadersFromContext(ctx),
	}
	err = ch.Publish(
		exchange, // exchange
		queue,    // routing key
		true,     // mandatory
		false,    // immediate
		msg,      // publishing
	)
	if err != nil {
		delete(r.rpcPending, correlationID)
		logging.LogError(err, "there was an error publishing an rpc message", "queue", queue)
		r.resetRPCClientLocked(ch, err)
		return err
	}
	timeoutChannel := time.After(RPC_TIMEOUT)
	for {
		select {
		case ntf := <-confirmChannel:
			if !ntf.Ack {
				err = errors.New("failed to deliver message, not ACK-ed by receiver")
				delete(r.rpcPending, correlationID)
				logging.LogError(err, "failed to deliver message to exchange/queue, notifyPublish", "queue", queue)
				r.resetRPCClientLocked(ch, err)
				return err
			}
			return nil
		case ret := <-notifyReturnChannel:
			if ret.CorrelationId != correlationID {
				continue
			}
			err = errors.New(getMeaningfulRabbitmqError(ret))
			delete(r.rpcPending, correlationID)
			r.resetRPCClientLocked(ch, err)
			return err
		case <-timeoutChannel:
			err = errors.New("message delivery confirmation timed out in SendRPCMessage")
			delete(r.rpcPending, correlationID)
			logging.LogError(err, "message delivery confirmation to exchange/queue timed out when sending", "queue", queue)
			r.resetRPCClientLocked(ch, err)
			return err
		}
	}
}

func (r *rabbitMQConnection) SendRPCMessage(exchange string, queue string, body []byte, exclusiveQueue bool, retryPolicies ...RPCRetryPolicy) ([]byte, error) {
	return r.SendRPCMessageWithContext(context.Background(), exchange, queue, body, exclusiveQueue, retryPolicies...)
}

func (r *rabbitMQConnection) SendRPCMessageWithContext(ctx context.Context, exchange string, queue string, body []byte, exclusiveQueue bool, retryPolicies ...RPCRetryPolicy) ([]byte, error) {
	var finalError error
	retryPolicy := getRPCRetryPolicy(retryPolicies)
	timeout := r.getRPCTimeout(retryPolicy)
	for attempt := 0; attempt < 3; attempt++ {
		correlationID := uuid.NewString()
		responseChannel := make(chan rpcResponse, 1)
		err := r.publishRPCMessage(ctx, exchange, queue, correlationID, body, exclusiveQueue, responseChannel)
		if err != nil {
			finalError = err
			time.Sleep(RPC_TIMEOUT)
			continue
		}
		//logging.LogDebug("Sent RPC message", "queue", queue)
		select {
		case response := <-responseChannel:
			if response.Err != nil {
				finalError = response.Err
				logging.LogError(response.Err, "rpc channel failed while waiting for response", "queue", queue)
				if retryPolicy != RPC_RETRY_POLICY_RETRY_ON_TIMEOUT {
					return nil, response.Err
				}
				time.Sleep(RPC_TIMEOUT)
				continue
			}
			//logging.LogDebug("Got RPC Reply", "queue", queue)
			return response.Body, nil
		case <-time.After(timeout):
			r.removePendingRPCResponse(correlationID)
			err = errors.New("rpc response timed out")
			finalError = err
			logging.LogError(err, "message delivery confirmation to exchange/queue timed out when receiving", "queue", queue)
			if retryPolicy != RPC_RETRY_POLICY_RETRY_ON_TIMEOUT {
				return nil, err
			}
			continue
		}
	}
	logging.LogError(finalError, "failed 3 times")
	return nil, finalError
}
func (r *rabbitMQConnection) ReceiveFromMythicDirectExchange(exchange string, queue string, routingKey string, handler QueueHandler, exclusiveQueue bool, wg *sync.WaitGroup) {
	// exchange is a direct exchange
	// queue is where the messages get sent to (local name)
	// routingKey is the specific direct topic we're interested in for the exchange
	// handler processes the messages we get on our queue
	for {
		conn, err := r.GetConnection()
		if err != nil {
			logging.LogError(err, "Failed to connect to rabbitmq", "retry_wait_time", RPC_TIMEOUT)
			time.Sleep(RPC_TIMEOUT)
			continue
		}
		ch, err := conn.Channel()
		if err != nil {
			logging.LogError(err, "Failed to open rabbitmq channel", "retry_wait_time", RPC_TIMEOUT)
			time.Sleep(RPC_TIMEOUT)
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
			logging.LogError(err, "Failed to declare exchange", "exchange", exchange, "exchange_type", "direct", "retry_wait_time", RPC_TIMEOUT)
			ch.Close()
			time.Sleep(RPC_TIMEOUT)
			continue
		}
		q, err := ch.QueueDeclare(
			queue,          // name, queue
			false,          // durable
			true,           // delete when unused
			exclusiveQueue, // exclusive
			false,          // no-wait
			nil,            // arguments
		)
		if err != nil {
			logging.LogError(err, "Failed to declare queue", "retry_wait_time", RPC_TIMEOUT, "queue", queue)
			ch.Close()
			time.Sleep(RPC_TIMEOUT)
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
			logging.LogError(err, "Failed to bind to queue to receive messages", "retry_wait_time", RPC_TIMEOUT, "queue", queue)
			ch.Close()
			time.Sleep(RPC_TIMEOUT)
			continue
		}
		msgs, err := ch.Consume(
			q.Name, // queue name
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no local
			false,  // no wait
			nil,    // args
		)
		if err != nil {
			logging.LogError(err, "Failed to start consuming on queue", "queue", q.Name)
			ch.Close()
			time.Sleep(RPC_TIMEOUT)
			continue
		}
		forever := make(chan bool)
		go func() {
			for d := range msgs {
				go handler(ContextFromHeaders(d.Headers), d.Body)
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
func (r *rabbitMQConnection) ReceiveFromRPCQueue(exchange string, queue string, routingKey string, handler RPCQueueHandler, exclusiveQueue bool, wg *sync.WaitGroup) {
	for {
		conn, err := r.GetConnection()
		if err != nil {
			logging.LogError(err, "Failed to connect to rabbitmq", "retry_wait_time", RPC_TIMEOUT, "queue", queue)
			time.Sleep(RPC_TIMEOUT)
			continue
		}
		ch, err := conn.Channel()
		if err != nil {
			logging.LogError(err, "Failed to open rabbitmq channel", "retry_wait_time", RPC_TIMEOUT, "queue", queue)
			time.Sleep(RPC_TIMEOUT)
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
			logging.LogError(err, "Failed to declare exchange", "exchange", exchange, "exchange_type", "direct", "retry_wait_time", RPC_TIMEOUT, "queue", queue)
			ch.Close()
			time.Sleep(RPC_TIMEOUT)
			continue
		}
		q, err := ch.QueueDeclare(
			queue,          // name, queue
			true,           // durable
			true,           // delete when unused
			exclusiveQueue, // exclusive
			false,          // no-wait
			nil,            // arguments
		)
		if err != nil {
			logging.LogError(err, "Failed to declare queue", "retry_wait_time", RPC_TIMEOUT, "queue", queue)
			ch.Close()
			time.Sleep(RPC_TIMEOUT)
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
			logging.LogError(err, "Failed to bind to queue to receive messages", "retry_wait_time", RPC_TIMEOUT, "queue", queue)
			ch.Close()
			time.Sleep(RPC_TIMEOUT)
			continue
		}
		forever := make(chan bool)
		go func() {
			for {
				if ch.IsClosed() {
					logging.LogError(nil, "channel is closed", "queue", q.Name)
					break
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
					time.Sleep(RETRY_CONNECT_DELAY)
					break
				}
				for d := range msgs {
					if ch.IsClosed() {
						logging.LogError(nil, "channel is closed", "queue", q.Name)
						forever <- true
						return
					}
					//logging.LogInfo("about to handle rpc msg", "queue", q.Name)
					responseMsg := handler(ContextFromHeaders(d.Headers), d.Body)
					//logging.LogInfo("finished handling rpc msg", "queue", q.Name)
					responseMsgJson, err := json.Marshal(responseMsg)
					if err != nil {
						logging.LogError(err, "Failed to generate JSON for rpc response", "queue", queue)
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
						logging.LogError(err, "Failed to send message", "queue", queue)
						continue
					}
					err = ch.Ack(d.DeliveryTag, false)
					if err != nil {
						logging.LogError(err, "Failed to Ack message", "queue", queue)
					}
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
	conn, err := r.GetConnection()
	if err != nil {
		return 0, err
	}
	ch, err := conn.Channel()
	if err != nil {
		logging.LogError(err, "Failed to open rabbitmq channel")
		return 0, err
	}
	defer ch.Close()
	err = ch.Confirm(false)
	if err != nil {
		logging.LogError(err, "Channel could not be put into confirm mode")
		return 0, err
	}
	err = ch.ExchangeDeclare(
		exchange, // exchange name
		kind,     // type of exchange, ex: topic, fanout, direct, etc
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logging.LogError(err, "Failed to declare exchange", "exchange", MYTHIC_TOPIC_EXCHANGE, "exchange_type", "topic", "retry_wait_time", RETRY_CONNECT_DELAY)
		return 0, err

	}
	q, err := ch.QueueDeclare(
		queue, // name, queue
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logging.LogError(err, "Unknown error (not 404 or 405) when checking for container existence")
		return 0, err
	}
	err = ch.QueueBind(
		q.Name,   // queue name
		queue,    // routing key
		exchange, // exchange name
		false,    // nowait
		nil,      // arguments
	)
	if err != nil {
		logging.LogError(err, "Failed to bind to queue to receive messages", "retry_wait_time", RETRY_CONNECT_DELAY)
		return 0, err
	}
	return uint(q.Consumers), nil
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
					go handler(ContextFromHeaders(d.Headers), d.Body)
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

func prepTaskArgs(ctx context.Context, command agentstructs.Command, taskMessage *agentstructs.PTTaskMessageAllData) error {
	args, err := agentstructs.GenerateArgsData(command.CommandParameters, *taskMessage)
	if err != nil {
		logging.LogError(err, "Failed to generate args data for create tasking")
		return errors.New(fmt.Sprintf("Failed to generate args data:\n%s", err.Error()))
	}
	taskMessage.Args = args
	if taskMessage.Task.IsInteractiveTask {
		// don't bother trying to parse arguments for interactive follow-ons, it doesn't make sense
		return nil
	}
	switch taskMessage.Args.GetTaskingLocation() {
	case "command_line":
		if command.TaskFunctionParseArgString != nil {
			// from scripting or if there are no parameters defined, this gets called
			err = command.TaskFunctionParseArgString(ctx, &taskMessage.Args, taskMessage.Args.GetCommandLine())
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
					if err := command.TaskFunctionParseArgString(ctx, &taskMessage.Args, taskMessage.Args.GetCommandLine()); err != nil {
						logging.LogError(err, "Failed to run ParseArgString function", "command Name", command.Name)
						return errors.New(fmt.Sprintf("Failed to run %s's ParseArgString function:\n%s", command.Name, err.Error()))
					} else {
						break
					}
				} else {
					logging.LogError(err, "Failed to parse arguments from parsed_cli into dictionary and no ParseArgString function defined, using raw command line")
				}

			} else if err := command.TaskFunctionParseArgDictionary(ctx, &taskMessage.Args, tempArgs); err != nil {
				logging.LogError(err, "Failed to run ParseArgDictionary function", "command Name", command.Name)
				return errors.New(fmt.Sprintf("Failed to run %s's ParseArgDictionary function:\n%s", command.Name, err.Error()))
			}
			// if no dictionary function, fall back to the string function if it exists
		} else if command.TaskFunctionParseArgString != nil {

			if err := command.TaskFunctionParseArgString(ctx, &taskMessage.Args, taskMessage.Args.GetCommandLine()); err != nil {
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
		newArg := arg.TypedArrayParseFunction(ctx, agentstructs.PTRPCTypedArrayParseFunctionMessage{
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
	stopResponse := C2RPCStopServer(context.Background(), c2structs.C2RPCStopServerMessage{
		Name: name,
	})
	if !stopResponse.Success {
		_, _ = SendMythicRPCC2UpdateStatus(context.Background(), MythicRPCC2UpdateStatusMessage{
			Error:                 stopResponse.Error,
			InternalServerRunning: stopResponse.InternalServerRunning,
			C2Profile:             name,
		})
		return
	}
	startResponse := C2RPCStartServer(context.Background(), c2structs.C2RPCStartServerMessage{
		Name: name,
	})
	_, _ = SendMythicRPCC2UpdateStatus(context.Background(), MythicRPCC2UpdateStatusMessage{
		Error:                 startResponse.Error,
		InternalServerRunning: startResponse.InternalServerRunning,
		C2Profile:             name,
	})
}

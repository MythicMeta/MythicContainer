package rabbitmq

import (
	"fmt"
	"github.com/MythicMeta/MythicContainer/grpc"
	"github.com/MythicMeta/MythicContainer/translationstructs"
	"os"
	"sync"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	c2structs "github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/loggingstructs"
	"github.com/MythicMeta/MythicContainer/utils"
	"github.com/MythicMeta/MythicContainer/webhookstructs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueHandler func([]byte)
type RPCQueueHandler func([]byte) interface{}
type RoutingKeyFunction func(string) string
type ContainerNameFunction func() string

type RPCQueueStruct struct {
	Exchange              string
	Queue                 string
	RoutingKeyFunction    RoutingKeyFunction
	ContainerNameFunction ContainerNameFunction
	Handler               RPCQueueHandler
}
type DirectQueueStruct struct {
	Exchange              string
	Queue                 string
	RoutingKeyFunction    RoutingKeyFunction
	ContainerNameFunction ContainerNameFunction
	Handler               QueueHandler
}

type rabbitMQConnection struct {
	conn             *amqp.Connection
	mutex            sync.RWMutex
	addListenerMutex sync.RWMutex
	RPCQueues        []RPCQueueStruct
	DirectQueues     []DirectQueueStruct
	needToResync     bool
}

var RabbitMQConnection rabbitMQConnection

const containerVersion = "v1.0.0-0.0.0"

func (r *rabbitMQConnection) AddRPCQueue(input RPCQueueStruct) {
	r.addListenerMutex.Lock()
	r.RPCQueues = append(r.RPCQueues, input)
	r.addListenerMutex.Unlock()
}
func (r *rabbitMQConnection) AddDirectQueue(input DirectQueueStruct) {
	r.addListenerMutex.Lock()
	r.DirectQueues = append(r.DirectQueues, input)
	r.addListenerMutex.Unlock()
}
func (r *rabbitMQConnection) startListeners(services []string) {
	// handle starting any queues that a developer isn't responsible for
	exclusiveQueue := true
	for _, rpcQueue := range r.RPCQueues {
		go RabbitMQConnection.ReceiveFromRPCQueue(
			rpcQueue.Exchange,
			rpcQueue.Queue,
			rpcQueue.RoutingKeyFunction(rpcQueue.ContainerNameFunction()),
			rpcQueue.Handler,
			exclusiveQueue)
	}
	for _, directQueue := range r.DirectQueues {
		go RabbitMQConnection.ReceiveFromMythicDirectExchange(
			directQueue.Exchange,
			directQueue.Queue,
			directQueue.RoutingKeyFunction(directQueue.ContainerNameFunction()),
			directQueue.Handler,
			exclusiveQueue)
	}
	// handle starting any queues that are necessary for the c2 profile
	if utils.StringSliceContains(services, "c2") {
		SyncAllC2Data(nil)
		for _, c2 := range c2structs.AllC2Data.GetAllNames() {
			// now that we're about to listen and sync, make sure all generic listeners are applied to all c2 profiles

			if c2structs.AllC2Data.Get(c2).GetC2Name() != "" {
				logging.LogInfo(fmt.Sprintf("Initializing RabbitMQ for C2 Service: %s\n", c2))
				for _, rpcQueue := range c2structs.AllC2Data.Get(c2).GetRPCMethods() {
					go RabbitMQConnection.ReceiveFromRPCQueue(
						MYTHIC_EXCHANGE,
						c2structs.AllC2Data.Get(c2).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						c2structs.AllC2Data.Get(c2).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						rpcQueue.RabbitmqProcessingFunction,
						exclusiveQueue,
					)
				}
				for _, rpcQueue := range c2structs.AllC2Data.Get("").GetRPCMethods() {
					go RabbitMQConnection.ReceiveFromRPCQueue(
						MYTHIC_EXCHANGE,
						c2structs.AllC2Data.Get(c2).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						c2structs.AllC2Data.Get(c2).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						rpcQueue.RabbitmqProcessingFunction,
						exclusiveQueue,
					)
				}
				for _, directQueue := range c2structs.AllC2Data.Get(c2).GetDirectMethods() {
					go RabbitMQConnection.ReceiveFromMythicDirectExchange(
						MYTHIC_EXCHANGE,
						c2structs.AllC2Data.Get(c2).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						c2structs.AllC2Data.Get(c2).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						directQueue.RabbitmqProcessingFunction,
						exclusiveQueue,
					)
				}
				for _, directQueue := range c2structs.AllC2Data.Get("").GetDirectMethods() {
					go RabbitMQConnection.ReceiveFromMythicDirectExchange(
						MYTHIC_EXCHANGE,
						c2structs.AllC2Data.Get(c2).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						c2structs.AllC2Data.Get(c2).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						directQueue.RabbitmqProcessingFunction,
						exclusiveQueue,
					)
				}
			} else {
				errorMessage := "Tasked C2 Container to start, but C2 agent name is empty.\n"
				errorMessage += "Did you initialize your functions module?\n"
				logging.LogError(nil, errorMessage)
				os.Exit(1)
			}
		}
	}
	// handle starting any queues that are necessary for the payload type
	if utils.StringSliceContains(services, "payload") {
		agentstructs.AllPayloadData.Get("").AddDirectMethod(agentstructs.RabbitmqDirectMethod{
			RabbitmqRoutingKey:         PAYLOAD_BUILD_ROUTING_KEY,
			RabbitmqProcessingFunction: WrapPayloadBuild,
		})
		SyncPayloadData(nil)
		for _, pt := range agentstructs.AllPayloadData.GetAllPayloadTypeNames() {

			if agentstructs.AllPayloadData.Get(pt).GetPayloadName() != "" {
				logging.LogInfo(fmt.Sprintf("Initializing RabbitMQ for Payload Service: %s\n", pt))
				for _, rpcQueue := range agentstructs.AllPayloadData.Get(pt).GetRPCMethods() {
					go RabbitMQConnection.ReceiveFromRPCQueue(
						MYTHIC_EXCHANGE,
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						rpcQueue.RabbitmqProcessingFunction,
						!exclusiveQueue,
					)
				}
				for _, rpcQueue := range agentstructs.AllPayloadData.Get("").GetRPCMethods() {
					go RabbitMQConnection.ReceiveFromRPCQueue(
						MYTHIC_EXCHANGE,
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						rpcQueue.RabbitmqProcessingFunction,
						!exclusiveQueue,
					)
				}
				for _, directQueue := range agentstructs.AllPayloadData.Get(pt).GetDirectMethods() {
					go RabbitMQConnection.ReceiveFromMythicDirectExchange(
						MYTHIC_EXCHANGE,
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						directQueue.RabbitmqProcessingFunction,
						!exclusiveQueue,
					)
				}
				for _, directQueue := range agentstructs.AllPayloadData.Get("").GetDirectMethods() {
					go RabbitMQConnection.ReceiveFromMythicDirectExchange(
						MYTHIC_EXCHANGE,
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						directQueue.RabbitmqProcessingFunction,
						!exclusiveQueue,
					)
				}
			} else {
				errorMessage := "Tasked Payload Container to start, but Payload agent name is empty.\n"
				errorMessage += "Did you initialize your functions module?\n"
				logging.LogError(nil, errorMessage)
				os.Exit(1)
			}
		}
	}
	// handle starting any queues that are necessary for the translation container
	if utils.StringSliceContains(services, "translation") {
		logging.LogInfo("Initializing RabbitMQ for Translation Services")
		SyncTranslationData(nil)
		for _, pt := range translationstructs.AllTranslationData.GetAllPayloadTypeNames() {

			if translationstructs.AllTranslationData.Get(pt).GetPayloadName() != "" {
				logging.LogInfo(fmt.Sprintf("Initializing RabbitMQ for Translation Service: %s\n", pt))
				go grpc.Initialize(translationstructs.AllTranslationData.Get(pt).GetPayloadName())
				for _, rpcQueue := range translationstructs.AllTranslationData.Get("").GetRPCMethods() {
					go RabbitMQConnection.ReceiveFromRPCQueue(
						MYTHIC_EXCHANGE,
						translationstructs.AllTranslationData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						translationstructs.AllTranslationData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						rpcQueue.RabbitmqProcessingFunction,
						!exclusiveQueue,
					)
				}
				/*
					for _, rpcQueue := range translationstructs.AllTranslationData.Get(pt).GetRPCMethods() {
						go RabbitMQConnection.ReceiveFromRPCQueue(
							MYTHIC_EXCHANGE,
							translationstructs.AllTranslationData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
							translationstructs.AllTranslationData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
							rpcQueue.RabbitmqProcessingFunction,
							!exclusiveQueue,
						)
					}

					for _, directQueue := range translationstructs.AllTranslationData.Get(pt).GetDirectMethods() {
						go RabbitMQConnection.ReceiveFromMythicDirectExchange(
							MYTHIC_EXCHANGE,
							translationstructs.AllTranslationData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
							translationstructs.AllTranslationData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
							directQueue.RabbitmqProcessingFunction,
							!exclusiveQueue,
						)
					}
					for _, directQueue := range translationstructs.AllTranslationData.Get("").GetDirectMethods() {
						go RabbitMQConnection.ReceiveFromMythicDirectExchange(
							MYTHIC_EXCHANGE,
							translationstructs.AllTranslationData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
							translationstructs.AllTranslationData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
							directQueue.RabbitmqProcessingFunction,
							!exclusiveQueue,
						)
					}

				*/
			} else {
				errorMessage := "Tasked Translation Container to start, but translation agent name is empty.\n"
				errorMessage += "Did you initialize your functions module?\n"
				logging.LogError(nil, errorMessage)
				os.Exit(1)
			}
		}
	}
	// handle starting any queues that are necessary for a logging container
	if utils.StringSliceContains(services, "logger") {
		logging.LogInfo("Initializing RabbitMQ for SIEM Logging Services")
		for _, directQueue := range loggingstructs.AllLoggingData.Get("").GetDirectMethods() {
			listenerExists := false
			for _, logger := range loggingstructs.AllLoggingData.GetAllNames() {
				loggingDef := loggingstructs.AllLoggingData.Get(logger).GetLoggingDefinition()
				switch directQueue.RabbitmqRoutingKey {
				case loggingstructs.LOG_TYPE_CALLBACK:
					if loggingDef.NewCallbackFunction != nil {
						listenerExists = true
					}
				case loggingstructs.LOG_TYPE_ARTIFACT:
					if loggingDef.NewArtifactFunction != nil {
						listenerExists = true
					}
				case loggingstructs.LOG_TYPE_CREDENTIAL:
					if loggingDef.NewCredentialFunction != nil {
						listenerExists = true
					}
				case loggingstructs.LOG_TYPE_KEYLOG:
					if loggingDef.NewKeylogFunction != nil {
						listenerExists = true
					}
				case loggingstructs.LOG_TYPE_FILE:
					if loggingDef.NewFileFunction != nil {
						listenerExists = true
					}
				case loggingstructs.LOG_TYPE_PAYLOAD:
					if loggingDef.NewPayloadFunction != nil {
						listenerExists = true
					}
				case loggingstructs.LOG_TYPE_TASK:
					if loggingDef.NewTaskFunction != nil {
						listenerExists = true
					}
				default:
				}
			}
			if listenerExists {
				go RabbitMQConnection.ReceiveFromMythicDirectTopicExchange(
					MYTHIC_TOPIC_EXCHANGE,
					loggingstructs.GetRoutingKeyFor(directQueue.RabbitmqRoutingKey),
					loggingstructs.GetRoutingKeyFor(directQueue.RabbitmqRoutingKey),
					directQueue.RabbitmqProcessingFunction,
					!exclusiveQueue,
				)
			}

		}
	}
	// handle starting any queues that are necessary for a webhook container
	if utils.StringSliceContains(services, "webhook") {
		logging.LogInfo("Initializing RabbitMQ for Webhook Services")
		for _, directQueue := range webhookstructs.AllWebhookData.Get("").GetDirectMethods() {
			listenerExists := false
			// only start listening for messages on queues if we have functions to process the messages
			for _, webhook := range webhookstructs.AllWebhookData.GetAllNames() {
				webhookDef := webhookstructs.AllWebhookData.Get(webhook).GetWebhookDefinition()
				switch directQueue.RabbitmqRoutingKey {
				case webhookstructs.WEBHOOK_TYPE_NEW_STARTUP:
					if webhookDef.NewStartupFunction != nil {
						listenerExists = true
					}
				case webhookstructs.WEBHOOK_TYPE_NEW_CALLBACK:
					if webhookDef.NewCallbackFunction != nil {
						listenerExists = true
					}
				case webhookstructs.WEBHOOK_TYPE_NEW_FEEDBACK:
					if webhookDef.NewFeedbackFunction != nil {
						listenerExists = true
					}
				default:
					logging.LogError(nil, "Unknown webhook type in rabbitmq initialize", "webhook type", directQueue.RabbitmqRoutingKey)
				}
			}
			if listenerExists {
				go RabbitMQConnection.ReceiveFromMythicDirectTopicExchange(
					MYTHIC_TOPIC_EXCHANGE,
					webhookstructs.GetRoutingKeyFor(directQueue.RabbitmqRoutingKey),
					webhookstructs.GetRoutingKeyFor(directQueue.RabbitmqRoutingKey),
					directQueue.RabbitmqProcessingFunction,
					!exclusiveQueue,
				)
			}
		}
	}
	logging.LogInfo("[+] All services initialized!")
}

func Initialize() {
	for {
		if _, err := RabbitMQConnection.GetConnection(); err == nil {
			logging.LogInfo("RabbitMQ Initialized")
			return
		}
		logging.LogInfo("Waiting for RabbitMQ...")
	}
}
func StartServices(services []string) {
	// define the exchange, mythic's queue name, which direct messages to get, and a function to handle messages for that queue
	// payload functionality
	RabbitMQConnection.startListeners(services)
}

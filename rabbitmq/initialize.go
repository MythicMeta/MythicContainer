package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/MythicMeta/MythicContainer/authstructs"
	"github.com/MythicMeta/MythicContainer/config"
	"github.com/MythicMeta/MythicContainer/eventingstructs"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"log"
	"os"
	"sync"

	"github.com/MythicMeta/MythicContainer/grpc"
	"github.com/MythicMeta/MythicContainer/loggingstructs"
	"github.com/MythicMeta/MythicContainer/translationstructs"
	"github.com/MythicMeta/MythicContainer/webhookstructs"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	c2structs "github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
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
var c2Mutex = sync.Mutex{}

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
	if config.MythicConfig.RabbitmqHost == "" {
		log.Fatalf("[-] Missing RABBITMQ_HOST environment variable point to rabbitmq server IP")
	}
	var wg sync.WaitGroup
	exclusiveQueue := true
	for _, rpcQueue := range r.RPCQueues {
		wg.Add(1)
		go RabbitMQConnection.ReceiveFromRPCQueue(
			rpcQueue.Exchange,
			rpcQueue.Queue,
			rpcQueue.RoutingKeyFunction(rpcQueue.ContainerNameFunction()),
			rpcQueue.Handler,
			exclusiveQueue,
			&wg)
	}
	for _, directQueue := range r.DirectQueues {
		wg.Add(1)
		go RabbitMQConnection.ReceiveFromMythicDirectExchange(
			directQueue.Exchange,
			directQueue.Queue,
			directQueue.RoutingKeyFunction(directQueue.ContainerNameFunction()),
			directQueue.Handler,
			exclusiveQueue,
			&wg)
	}
	wg.Wait()
	// handle starting any queues that are necessary for a logging container
	if helpers.StringSliceContains(services, "logger") {
		logging.LogInfo("Initializing RabbitMQ for SIEM Logging Services")
		for _, logger := range loggingstructs.AllLoggingData.GetAllNames() {
			// get our resync routing and register it
			loggingstructs.AllLoggingData.Get(logger).SetName(logger)
			for _, rpcQueue := range loggingstructs.AllLoggingData.Get("").GetRPCMethods() {
				go RabbitMQConnection.ReceiveFromRPCQueue(
					MYTHIC_EXCHANGE,
					loggingstructs.AllLoggingData.Get(logger).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
					loggingstructs.AllLoggingData.Get(logger).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
					rpcQueue.RabbitmqProcessingFunction,
					exclusiveQueue,
					nil,
				)
			}
			loggingDef := loggingstructs.AllLoggingData.Get(logger).GetLoggingDefinition()
			subscriptions := []string{}
			for _, directQueue := range loggingstructs.AllLoggingData.Get("").GetDirectMethods() {
				listenerExists := false
				switch directQueue.RabbitmqRoutingKey {
				case loggingstructs.LOG_TYPE_CALLBACK:
					if loggingDef.NewCallbackFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, loggingstructs.LOG_TYPE_CALLBACK)
					}
				case loggingstructs.LOG_TYPE_ARTIFACT:
					if loggingDef.NewArtifactFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, loggingstructs.LOG_TYPE_ARTIFACT)
					}
				case loggingstructs.LOG_TYPE_CREDENTIAL:
					if loggingDef.NewCredentialFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, loggingstructs.LOG_TYPE_CREDENTIAL)
					}
				case loggingstructs.LOG_TYPE_KEYLOG:
					if loggingDef.NewKeylogFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, loggingstructs.LOG_TYPE_KEYLOG)
					}
				case loggingstructs.LOG_TYPE_FILE:
					if loggingDef.NewFileFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, loggingstructs.LOG_TYPE_FILE)
					}
				case loggingstructs.LOG_TYPE_PAYLOAD:
					if loggingDef.NewPayloadFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, loggingstructs.LOG_TYPE_PAYLOAD)
					}
				case loggingstructs.LOG_TYPE_TASK:
					if loggingDef.NewTaskFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, loggingstructs.LOG_TYPE_TASK)
					}
				case loggingstructs.LOG_TYPE_RESPONSE:
					if loggingDef.NewResponseFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, loggingstructs.LOG_TYPE_RESPONSE)
					}
				default:
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
			loggingstructs.AllLoggingData.Get(logger).SetSubscriptions(subscriptions)
			SyncConsumingContainerData(logger, "logging")
		}
	}
	// handle starting any queues that are necessary for a webhook container
	if helpers.StringSliceContains(services, "webhook") {
		logging.LogInfo("Initializing RabbitMQ for Webhook Services")
		for _, webhook := range webhookstructs.AllWebhookData.GetAllNames() {
			webhookstructs.AllWebhookData.Get(webhook).SetName(webhook)
			// get our resync routing and register it
			for _, rpcQueue := range webhookstructs.AllWebhookData.Get("").GetRPCMethods() {
				go RabbitMQConnection.ReceiveFromRPCQueue(
					MYTHIC_EXCHANGE,
					webhookstructs.AllWebhookData.Get(webhook).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
					webhookstructs.AllWebhookData.Get(webhook).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
					rpcQueue.RabbitmqProcessingFunction,
					exclusiveQueue,
					nil,
				)
			}
			webhookDef := webhookstructs.AllWebhookData.Get(webhook).GetWebhookDefinition()
			subscriptions := []string{}
			for _, directQueue := range webhookstructs.AllWebhookData.Get("").GetDirectMethods() {
				listenerExists := false
				// only start listening for messages on queues if we have functions to process the messages
				switch directQueue.RabbitmqRoutingKey {
				case webhookstructs.WEBHOOK_TYPE_NEW_STARTUP:
					if webhookDef.NewStartupFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, webhookstructs.WEBHOOK_TYPE_NEW_STARTUP)
					}
				case webhookstructs.WEBHOOK_TYPE_NEW_CALLBACK:
					if webhookDef.NewCallbackFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, webhookstructs.WEBHOOK_TYPE_NEW_CALLBACK)
					}
				case webhookstructs.WEBHOOK_TYPE_NEW_FEEDBACK:
					if webhookDef.NewFeedbackFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, webhookstructs.WEBHOOK_TYPE_NEW_FEEDBACK)
					}
				case webhookstructs.WEBHOOK_TYPE_NEW_ALERT:
					if webhookDef.NewAlertFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, webhookstructs.WEBHOOK_TYPE_NEW_ALERT)
					}
				case webhookstructs.WEBHOOK_TYPE_NEW_CUSTOM:
					if webhookDef.NewCustomFunction != nil {
						listenerExists = true
						subscriptions = append(subscriptions, webhookstructs.WEBHOOK_TYPE_NEW_CUSTOM)
					}
				default:
					logging.LogError(nil, "Unknown webhook type in rabbitmq initialize", "webhook type", directQueue.RabbitmqRoutingKey)
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
			webhookstructs.AllWebhookData.Get(webhook).SetSubscriptions(subscriptions)
			SyncConsumingContainerData(webhook, "webhook")
		}
	}
	// handle starting any queues that are necessary for the c2 profile
	if helpers.StringSliceContains(services, "c2") {
		//SyncAllC2Data(nil)
		for _, c2 := range c2structs.AllC2Data.GetAllNames() {
			// now that we're about to listen and sync, make sure all generic listeners are applied to all c2 profiles
			if c2structs.AllC2Data.Get(c2).GetC2Name() != "" {
				logging.LogInfo(fmt.Sprintf("Initializing RabbitMQ for C2 Service: %s\n", c2))
				for _, rpcQueue := range c2structs.AllC2Data.Get(c2).GetRPCMethods() {
					wg.Add(1)
					go RabbitMQConnection.ReceiveFromRPCQueue(
						MYTHIC_EXCHANGE,
						c2structs.AllC2Data.Get(c2).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						c2structs.AllC2Data.Get(c2).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						rpcQueue.RabbitmqProcessingFunction,
						exclusiveQueue,
						&wg,
					)
				}
				for _, rpcQueue := range c2structs.AllC2Data.Get("").GetRPCMethods() {
					wg.Add(1)
					go RabbitMQConnection.ReceiveFromRPCQueue(
						MYTHIC_EXCHANGE,
						c2structs.AllC2Data.Get(c2).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						c2structs.AllC2Data.Get(c2).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						rpcQueue.RabbitmqProcessingFunction,
						exclusiveQueue,
						&wg,
					)
				}
				for _, directQueue := range c2structs.AllC2Data.Get(c2).GetDirectMethods() {
					wg.Add(1)
					go RabbitMQConnection.ReceiveFromMythicDirectExchange(
						MYTHIC_EXCHANGE,
						c2structs.AllC2Data.Get(c2).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						c2structs.AllC2Data.Get(c2).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						directQueue.RabbitmqProcessingFunction,
						exclusiveQueue,
						&wg,
					)
				}
				for _, directQueue := range c2structs.AllC2Data.Get("").GetDirectMethods() {
					wg.Add(1)
					go RabbitMQConnection.ReceiveFromMythicDirectExchange(
						MYTHIC_EXCHANGE,
						c2structs.AllC2Data.Get(c2).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						c2structs.AllC2Data.Get(c2).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						directQueue.RabbitmqProcessingFunction,
						exclusiveQueue,
						&wg,
					)
				}
				wg.Wait()
				SyncAllC2Data(&c2)
			} else {
				errorMessage := "Tasked C2 Container to start, but C2 agent name is empty.\n"
				errorMessage += "Did you initialize your functions module?\n"
				logging.LogError(nil, errorMessage)
				os.Exit(1)
			}
		}
	}
	// handle starting any queues that are necessary for the payload type
	if helpers.StringSliceContains(services, "payload") {
		agentstructs.AllPayloadData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
			RabbitmqRoutingKey:         PAYLOAD_BUILD_ROUTING_KEY,
			RabbitmqProcessingFunction: WrapPayloadBuild,
		})
		var PTwg sync.WaitGroup
		for _, pt := range agentstructs.AllPayloadData.GetAllPayloadTypeNames() {

			if agentstructs.AllPayloadData.Get(pt).GetPayloadName() != "" {
				logging.LogInfo(fmt.Sprintf("Initializing RabbitMQ for Payload Service: %s\n", pt))
				for _, rpcQueue := range agentstructs.AllPayloadData.Get(pt).GetRPCMethods() {
					PTwg.Add(1)
					go RabbitMQConnection.ReceiveFromRPCQueue(
						MYTHIC_EXCHANGE,
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						rpcQueue.RabbitmqProcessingFunction,
						!exclusiveQueue,
						&PTwg,
					)
				}
				for _, rpcQueue := range agentstructs.AllPayloadData.Get("").GetRPCMethods() {
					PTwg.Add(1)
					go RabbitMQConnection.ReceiveFromRPCQueue(
						MYTHIC_EXCHANGE,
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
						rpcQueue.RabbitmqProcessingFunction,
						!exclusiveQueue,
						&PTwg,
					)
				}
				for _, directQueue := range agentstructs.AllPayloadData.Get(pt).GetDirectMethods() {
					PTwg.Add(1)
					go RabbitMQConnection.ReceiveFromMythicDirectExchange(
						MYTHIC_EXCHANGE,
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						directQueue.RabbitmqProcessingFunction,
						!exclusiveQueue,
						&PTwg,
					)
				}
				for _, directQueue := range agentstructs.AllPayloadData.Get("").GetDirectMethods() {
					PTwg.Add(1)
					go RabbitMQConnection.ReceiveFromMythicDirectExchange(
						MYTHIC_EXCHANGE,
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						agentstructs.AllPayloadData.Get(pt).GetRoutingKey(directQueue.RabbitmqRoutingKey),
						directQueue.RabbitmqProcessingFunction,
						!exclusiveQueue,
						&PTwg,
					)
				}
			} else {
				errorMessage := "Tasked Payload Container to start, but Payload agent name is empty.\n"
				errorMessage += "Did you initialize your functions module?\n"
				logging.LogError(nil, errorMessage)
				os.Exit(1)
			}
		}
		wg.Wait()
		SyncPayloadData(nil, false)
	}
	// handle starting any queues that are necessary for the translation container
	if helpers.StringSliceContains(services, "translation") {
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
						nil,
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
	// handle starting any queues that are necessary for eventing containers
	if helpers.StringSliceContains(services, "eventing") {
		logging.LogInfo("Initializing RabbitMQ for Eventing Services")
		for _, eventer := range eventingstructs.AllEventingData.GetAllNames() {
			// get our resync routing and register it
			eventingstructs.AllEventingData.Get(eventer).SetName(eventer)
			for _, rpcQueue := range eventingstructs.AllEventingData.Get("").GetRPCMethods() {
				go RabbitMQConnection.ReceiveFromRPCQueue(
					MYTHIC_EXCHANGE,
					eventingstructs.AllEventingData.Get(eventer).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
					eventingstructs.AllEventingData.Get(eventer).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
					rpcQueue.RabbitmqProcessingFunction,
					exclusiveQueue,
					nil,
				)
			}
			eventingDef := eventingstructs.AllEventingData.Get(eventer).GetEventingDefinition()
			subscriptions := []string{}
			for i, _ := range eventingDef.CustomFunctions {
				subBytes, err := json.Marshal(eventingDef.CustomFunctions[i])
				if err != nil {
					logging.LogError(err, "Failed to marshal eventing function definition")
				} else {
					subscriptions = append(subscriptions, string(subBytes))
				}
			}
			for i, _ := range eventingDef.ConditionalChecks {
				subBytes, err := json.Marshal(eventingDef.ConditionalChecks[i])
				if err != nil {
					logging.LogError(err, "Failed to marshal eventing function definition")
				} else {
					subscriptions = append(subscriptions, string(subBytes))
				}
			}
			if eventingDef.TaskInterceptFunction != nil {
				subBytes, err := json.Marshal(map[string]string{
					"name":        "task_intercept",
					"description": "Intercept Task execution before it gets to an agent and potentially block it",
				})
				if err != nil {
					logging.LogError(err, "Failed to marshal eventing function definition")
				} else {
					subscriptions = append(subscriptions, string(subBytes))
				}
			}
			if eventingDef.ResponseInterceptFunction != nil {
				subBytes, err := json.Marshal(map[string]string{
					"name":        "response_intercept",
					"description": "Intercept User Output Responses they get sent to the user",
				})
				if err != nil {
					logging.LogError(err, "Failed to marshal eventing function definition")
				} else {
					subscriptions = append(subscriptions, string(subBytes))
				}
			}
			for _, directQueue := range eventingstructs.AllEventingData.Get("").GetDirectMethods() {
				go RabbitMQConnection.ReceiveFromMythicDirectExchange(
					MYTHIC_EXCHANGE,
					eventingstructs.AllEventingData.Get(eventer).GetRoutingKey(directQueue.RabbitmqRoutingKey),
					eventingstructs.AllEventingData.Get(eventer).GetRoutingKey(directQueue.RabbitmqRoutingKey),
					directQueue.RabbitmqProcessingFunction,
					exclusiveQueue,
					nil,
				)
			}
			eventingstructs.AllEventingData.Get(eventer).SetSubscriptions(subscriptions)
			SyncConsumingContainerData(eventer, "eventing")
		}
	}
	// handle starting any queues that are necessary for auth containers
	if helpers.StringSliceContains(services, "auth") {
		logging.LogInfo("Initializing RabbitMQ for Auth Services")
		for _, eventer := range authstructs.AllAuthData.GetAllNames() {
			// get our resync routing and register it
			authstructs.AllAuthData.Get(eventer).SetName(eventer)
			for _, rpcQueue := range authstructs.AllAuthData.Get("").GetRPCMethods() {
				go RabbitMQConnection.ReceiveFromRPCQueue(
					MYTHIC_EXCHANGE,
					authstructs.AllAuthData.Get(eventer).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
					authstructs.AllAuthData.Get(eventer).GetRoutingKey(rpcQueue.RabbitmqRoutingKey),
					rpcQueue.RabbitmqProcessingFunction,
					exclusiveQueue,
					nil,
				)
			}
			for _, directQueue := range authstructs.AllAuthData.Get("").GetDirectMethods() {
				go RabbitMQConnection.ReceiveFromMythicDirectExchange(
					MYTHIC_EXCHANGE,
					authstructs.AllAuthData.Get(eventer).GetRoutingKey(directQueue.RabbitmqRoutingKey),
					authstructs.AllAuthData.Get(eventer).GetRoutingKey(directQueue.RabbitmqRoutingKey),
					directQueue.RabbitmqProcessingFunction,
					exclusiveQueue,
					nil,
				)
			}
			authDef := authstructs.AllAuthData.Get(eventer).GetAuthDefinition()
			subscriptions := []string{}
			for _, sub := range authDef.IDPServices {
				subBytes, err := json.Marshal(map[string]string{
					"name": sub,
					"type": "idp",
				})
				if err != nil {
					logging.LogError(err, "Failed to marshal auth IDP Service definition")
				} else {
					subscriptions = append(subscriptions, string(subBytes))
				}
			}
			for _, sub := range authDef.NonIDPServices {
				subBytes, err := json.Marshal(map[string]string{
					"name": sub,
					"type": "nonidp",
				})
				if err != nil {
					logging.LogError(err, "Failed to marshal auth NonIDP Service definition")
				} else {
					subscriptions = append(subscriptions, string(subBytes))
				}
			}
			authstructs.AllAuthData.Get(eventer).SetSubscriptions(subscriptions)
			SyncConsumingContainerData(eventer, "auth")
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
	logging.LogInfo("Starting Services", "containerVersion", containerVersion)
	logging.LogInfo(containerVersionMessage)
	RabbitMQConnection.startListeners(services)
	logging.LogInfo("Successfully Started", "containerVersion", containerVersion)
}

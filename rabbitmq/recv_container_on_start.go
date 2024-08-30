package rabbitmq

import (
	"encoding/json"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/authstructs"
	c2structs "github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/eventingstructs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/loggingstructs"
	"github.com/MythicMeta/MythicContainer/translationstructs"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"github.com/MythicMeta/MythicContainer/webhookstructs"
)

// Register this method with rabbitmq so it can be called
func init() {
	eventingstructs.AllEventingData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CONTAINER_ON_START,
		RabbitmqProcessingFunction: processOnEventingStart,
	})
	agentstructs.AllPayloadData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CONTAINER_ON_START,
		RabbitmqProcessingFunction: processOnEventingStart,
	})
	c2structs.AllC2Data.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CONTAINER_ON_START,
		RabbitmqProcessingFunction: processOnEventingStart,
	})
	loggingstructs.AllLoggingData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CONTAINER_ON_START,
		RabbitmqProcessingFunction: processOnEventingStart,
	})
	translationstructs.AllTranslationData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CONTAINER_ON_START,
		RabbitmqProcessingFunction: processOnEventingStart,
	})
	webhookstructs.AllWebhookData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CONTAINER_ON_START,
		RabbitmqProcessingFunction: processOnEventingStart,
	})
	authstructs.AllAuthData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CONTAINER_ON_START,
		RabbitmqProcessingFunction: processOnEventingStart,
	})
}
func processOnEventingStart(input []byte) {
	inputStruct := sharedStructs.ContainerOnStartMessage{}
	err := json.Unmarshal(input, &inputStruct)
	if err != nil {
		logging.LogError(err, "Failed to process new callback webhook message")
		return
	}
	for _, containerName := range agentstructs.AllPayloadData.GetAllPayloadTypeNames() {
		if agentstructs.AllPayloadData.Get(containerName).GetPayloadDefinition().Name == inputStruct.ContainerName {
			if agentstructs.AllPayloadData.Get(containerName).GetPayloadDefinition().OnContainerStartFunction != nil {
				go processContainerOnStart(agentstructs.AllPayloadData.Get(containerName).GetPayloadDefinition().OnContainerStartFunction,
					inputStruct)
				return
			}
		}
	}
	for _, containerName := range c2structs.AllC2Data.GetAllNames() {
		if c2structs.AllC2Data.Get(containerName).GetC2Definition().Name == inputStruct.ContainerName {
			if c2structs.AllC2Data.Get(containerName).GetC2Definition().OnContainerStartFunction != nil {
				go processContainerOnStart(c2structs.AllC2Data.Get(containerName).GetC2Definition().OnContainerStartFunction,
					inputStruct)
				return
			}
		}
	}
	for _, containerName := range loggingstructs.AllLoggingData.GetAllNames() {
		if loggingstructs.AllLoggingData.Get(containerName).GetLoggingDefinition().Name == inputStruct.ContainerName {
			if loggingstructs.AllLoggingData.Get(containerName).GetLoggingDefinition().OnContainerStartFunction != nil {
				go processContainerOnStart(loggingstructs.AllLoggingData.Get(containerName).GetLoggingDefinition().OnContainerStartFunction,
					inputStruct)
				return
			}
		}
	}
	for _, containerName := range translationstructs.AllTranslationData.GetAllPayloadTypeNames() {
		if translationstructs.AllTranslationData.Get(containerName).GetPayloadDefinition().Name == inputStruct.ContainerName {
			if translationstructs.AllTranslationData.Get(containerName).GetPayloadDefinition().OnContainerStartFunction != nil {
				go processContainerOnStart(translationstructs.AllTranslationData.Get(containerName).GetPayloadDefinition().OnContainerStartFunction,
					inputStruct)
				return
			}
		}
	}
	for _, containerName := range webhookstructs.AllWebhookData.GetAllNames() {
		if webhookstructs.AllWebhookData.Get(containerName).GetWebhookDefinition().Name == inputStruct.ContainerName {
			if webhookstructs.AllWebhookData.Get(containerName).GetWebhookDefinition().OnContainerStartFunction != nil {
				go processContainerOnStart(webhookstructs.AllWebhookData.Get(containerName).GetWebhookDefinition().OnContainerStartFunction,
					inputStruct)
				return
			}
		}
	}
	for _, containerName := range eventingstructs.AllEventingData.GetAllNames() {
		if eventingstructs.AllEventingData.Get(containerName).GetEventingDefinition().Name == inputStruct.ContainerName {
			if eventingstructs.AllEventingData.Get(containerName).GetEventingDefinition().OnContainerStartFunction != nil {
				go processContainerOnStart(eventingstructs.AllEventingData.Get(containerName).GetEventingDefinition().OnContainerStartFunction,
					inputStruct)
				return
			}
		}
	}
	for _, containerName := range authstructs.AllAuthData.GetAllNames() {
		if authstructs.AllAuthData.Get(containerName).GetAuthDefinition().Name == inputStruct.ContainerName {
			if authstructs.AllAuthData.Get(containerName).GetAuthDefinition().OnContainerStartFunction != nil {
				go processContainerOnStart(authstructs.AllAuthData.Get(containerName).GetAuthDefinition().OnContainerStartFunction,
					inputStruct)
				return
			}
		}
	}
}
func processContainerOnStart(inputFunc func(sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse,
	incomingMessage sharedStructs.ContainerOnStartMessage) {
	response := inputFunc(incomingMessage)
	response.ContainerName = incomingMessage.ContainerName
	err := RabbitMQConnection.SendStructMessage(
		MYTHIC_EXCHANGE,
		CONTAINER_ON_START_RESPONSE,
		"",
		response,
		false,
	)
	if err != nil {
		logging.LogError(err, "Failed to send response back to Mythic")
	}
}

package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/MythicMeta/MythicContainer/eventingstructs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

// Register this method with rabbitmq so it can be called
func init() {
	eventingstructs.AllEventingData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         EVENTING_TASK_INTERCEPT,
		RabbitmqProcessingFunction: processTaskInterceptEventingFunction,
	})
}
func processTaskInterceptEventingFunction(input []byte) {
	inputStruct := eventingstructs.TaskInterceptMessage{}
	err := json.Unmarshal(input, &inputStruct)
	if err != nil {
		logging.LogError(err, "Failed to process new callback webhook message")
		return
	}
	for _, eventing := range eventingstructs.AllEventingData.GetAllNames() {
		if eventingstructs.AllEventingData.Get(eventing).GetEventingDefinition().Name == inputStruct.ContainerName {
			if eventingstructs.AllEventingData.Get(eventing).GetEventingDefinition().TaskInterceptFunction != nil {
				go func(incomingMessage eventingstructs.TaskInterceptMessage) {
					response := eventingstructs.AllEventingData.Get(eventing).GetEventingDefinition().TaskInterceptFunction(incomingMessage)
					response.EventStepInstanceID = incomingMessage.EventStepInstanceID
					response.TaskID = incomingMessage.TaskID
					if response.BlockTask && response.BypassRole == "" {
						response.BypassRole = eventingstructs.OPSEC_ROLE_OPERATOR
					}
					err = RabbitMQConnection.SendStructMessage(
						MYTHIC_EXCHANGE,
						EVENTING_TASK_INTERCEPT_RESPONSE,
						"",
						response,
						false,
					)
					if err != nil {
						logging.LogError(err, "Failed to send payload response back to Mythic")
					}
				}(inputStruct)
				return
			}
			logging.LogError(nil, fmt.Sprintf("Found container name, %s, but missing task intercept function",
				inputStruct.ContainerName))
			err = RabbitMQConnection.SendStructMessage(
				MYTHIC_EXCHANGE,
				EVENTING_TASK_INTERCEPT_RESPONSE,
				"",
				eventingstructs.TaskInterceptMessageResponse{
					Success:             false,
					EventStepInstanceID: inputStruct.EventStepInstanceID,
					StdErr: fmt.Sprintf("Found container name, %s, but missing task intercept function",
						inputStruct.ContainerName),
				},
				false,
			)
		}
	}
}

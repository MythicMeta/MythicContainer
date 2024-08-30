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
		RabbitmqRoutingKey:         EVENTING_CUSTOM_FUNCTION,
		RabbitmqProcessingFunction: processNewCustomEventingFunction,
	})
}
func processNewCustomEventingFunction(input []byte) {
	inputStruct := eventingstructs.NewCustomEventingMessage{}
	err := json.Unmarshal(input, &inputStruct)
	if err != nil {
		logging.LogError(err, "Failed to process new callback webhook message")
		return
	}
	for _, eventing := range eventingstructs.AllEventingData.GetAllNames() {
		if eventingstructs.AllEventingData.Get(eventing).GetEventingDefinition().Name == inputStruct.ContainerName {
			for i, _ := range eventingstructs.AllEventingData.Get(eventing).GetEventingDefinition().CustomFunctions {
				if eventingstructs.AllEventingData.Get(eventing).GetEventingDefinition().CustomFunctions[i].Name == inputStruct.FunctionName {
					go func(incomingMessage eventingstructs.NewCustomEventingMessage) {
						response := eventingstructs.AllEventingData.Get(eventing).GetEventingDefinition().CustomFunctions[i].Function(incomingMessage)
						response.EventStepInstanceID = incomingMessage.EventStepInstanceID
						err = RabbitMQConnection.SendStructMessage(
							MYTHIC_EXCHANGE,
							EVENTING_CUSTOM_FUNCTION_RESPONSE,
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
			}
			logging.LogError(nil, fmt.Sprintf("Found container name, %s, but missing function, %s",
				inputStruct.ContainerName, inputStruct.FunctionName))
			err = RabbitMQConnection.SendStructMessage(
				MYTHIC_EXCHANGE,
				EVENTING_CUSTOM_FUNCTION_RESPONSE,
				"",
				eventingstructs.NewCustomEventingMessageResponse{
					Success:             false,
					EventStepInstanceID: inputStruct.EventStepInstanceID,
					StdErr: fmt.Sprintf("Found container name, %s, but missing function, %s",
						inputStruct.ContainerName, inputStruct.FunctionName),
				},
				false,
			)
		}
	}
}

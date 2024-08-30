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
		RabbitmqRoutingKey:         EVENTING_RESPONSE_INTERCEPT,
		RabbitmqProcessingFunction: processResponseInterceptEventingFunction,
	})
}
func processResponseInterceptEventingFunction(input []byte) {
	inputStruct := eventingstructs.ResponseInterceptMessage{}
	err := json.Unmarshal(input, &inputStruct)
	if err != nil {
		logging.LogError(err, "Failed to process new callback webhook message")
		return
	}
	for _, eventing := range eventingstructs.AllEventingData.GetAllNames() {
		if eventingstructs.AllEventingData.Get(eventing).GetEventingDefinition().Name == inputStruct.ContainerName {
			if eventingstructs.AllEventingData.Get(eventing).GetEventingDefinition().ResponseInterceptFunction != nil {
				go func(incomingMessage eventingstructs.ResponseInterceptMessage) {
					response := eventingstructs.AllEventingData.Get(eventing).GetEventingDefinition().ResponseInterceptFunction(incomingMessage)
					response.EventStepInstanceID = incomingMessage.EventStepInstanceID
					response.ResponseID = incomingMessage.ResponseID
					err = RabbitMQConnection.SendStructMessage(
						MYTHIC_EXCHANGE,
						EVENTING_RESPONSE_INTERCEPT_RESPONSE,
						"",
						response,
						false,
					)
					if err != nil {
						logging.LogError(err, "Failed to send response intercept response back to Mythic")
					}
				}(inputStruct)
				return
			}
			logging.LogError(nil, fmt.Sprintf("Found container name, %s, but missing response intercept function",
				inputStruct.ContainerName))
			err = RabbitMQConnection.SendStructMessage(
				MYTHIC_EXCHANGE,
				EVENTING_RESPONSE_INTERCEPT_RESPONSE,
				"",
				eventingstructs.ResponseInterceptMessageResponse{
					Success:             false,
					EventStepInstanceID: inputStruct.EventStepInstanceID,
					ResponseID:          inputStruct.ResponseID,
					StdErr: fmt.Sprintf("Found container name, %s, but missing task intercept function",
						inputStruct.ContainerName),
				},
				false,
			)
		}
	}
}

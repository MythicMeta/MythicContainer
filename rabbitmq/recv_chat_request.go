package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MythicMeta/MythicContainer/chatstructs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

func init() {
	chatstructs.AllChatData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CHAT_REQUEST,
		RabbitmqProcessingFunction: processChatRequest,
	})
}

func processChatRequest(ctx context.Context, input []byte) {
	inputStruct := chatstructs.ChatRequestMessage{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process chat request message")
		return
	}
	for _, chat := range chatstructs.AllChatData.GetAllNames() {
		chatDefinition := chatstructs.AllChatData.Get(chat).GetChatDefinition()
		if chatDefinition.Name == inputStruct.ContainerName {
			if chatDefinition.ChatFunction != nil {
				go chatDefinition.ChatFunction(ctx, inputStruct)
				return
			}
			sendChatResponse(ctx, chatstructs.ChatResponseMessage{
				OperationID:       inputStruct.OperationID,
				RequestID:         inputStruct.RequestID,
				ResponseMessageID: inputStruct.ResponseMessageID,
				Status:            "error",
				Error:             fmt.Sprintf("%s does not implement a chat function", inputStruct.ContainerName),
			})
			return
		}
	}
	sendChatResponse(ctx, chatstructs.ChatResponseMessage{
		OperationID:       inputStruct.OperationID,
		RequestID:         inputStruct.RequestID,
		ResponseMessageID: inputStruct.ResponseMessageID,
		Status:            "error",
		Error:             fmt.Sprintf("Failed to find chat service %s", inputStruct.ContainerName),
	})
}

func sendChatResponse(ctx context.Context, response chatstructs.ChatResponseMessage) {
	if err := RabbitMQConnection.SendStructMessageWithContext(ctx, MYTHIC_EXCHANGE, CHAT_RESPONSE_ROUTING_KEY, "", response, false); err != nil {
		logging.LogError(err, "Failed to send chat response back to Mythic")
	}
}

func SendChatResponse(ctx context.Context, response chatstructs.ChatResponseMessage) {
	sendChatResponse(ctx, response)
}

package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/MythicMeta/MythicContainer/chatstructs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

var activeChatRequests = struct {
	sync.Mutex
	cancelByRequestID map[int]context.CancelFunc
	cancelledRequests map[int]struct{}
}{
	cancelByRequestID: map[int]context.CancelFunc{},
	cancelledRequests: map[int]struct{}{},
}

func init() {
	chatstructs.AllChatData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CHAT_REQUEST,
		RabbitmqProcessingFunction: processChatRequest,
	})
	chatstructs.AllChatData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CHAT_CANCEL,
		RabbitmqProcessingFunction: processChatCancel,
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
				requestCtx, cancel := context.WithCancel(ctx)
				if !canRegisterChatCancel(inputStruct.RequestID, cancel) {
					cancel()
					unregisterChatCancel(inputStruct.RequestID)
					sendChatResponse(ctx, chatstructs.ChatResponseMessage{
						OperationID:       inputStruct.OperationID,
						RequestID:         inputStruct.RequestID,
						ResponseMessageID: inputStruct.ResponseMessageID,
						Status:            "cancelled",
						Error:             "Cancelled by operator",
					})
					return
				}
				go func() {
					defer unregisterChatCancel(inputStruct.RequestID)
					chatDefinition.ChatFunction(requestCtx, inputStruct)
				}()
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

func processChatCancel(ctx context.Context, input []byte) {
	inputStruct := chatstructs.ChatCancelMessage{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process chat cancellation message")
		return
	}
	if inputStruct.RequestID <= 0 {
		return
	}
	activeChatRequests.Lock()
	cancel := activeChatRequests.cancelByRequestID[inputStruct.RequestID]
	if cancel == nil {
		activeChatRequests.cancelledRequests[inputStruct.RequestID] = struct{}{}
	}
	activeChatRequests.Unlock()
	if cancel != nil {
		cancel()
	}
}

func canRegisterChatCancel(requestID int, cancel context.CancelFunc) bool {
	if requestID <= 0 {
		return false
	}
	activeChatRequests.Lock()
	defer activeChatRequests.Unlock()
	if _, ok := activeChatRequests.cancelledRequests[requestID]; ok {
		delete(activeChatRequests.cancelledRequests, requestID)
		return false
	}
	activeChatRequests.cancelByRequestID[requestID] = cancel
	return true
}

func unregisterChatCancel(requestID int) {
	if requestID <= 0 {
		return
	}
	activeChatRequests.Lock()
	delete(activeChatRequests.cancelByRequestID, requestID)
	activeChatRequests.Unlock()
}

func sendChatResponse(ctx context.Context, response chatstructs.ChatResponseMessage) {
	if err := RabbitMQConnection.SendStructMessageWithContext(ctx, MYTHIC_EXCHANGE, CHAT_RESPONSE_ROUTING_KEY, "", response, false); err != nil {
		logging.LogError(err, "Failed to send chat response back to Mythic")
	}
}

func SendChatResponse(ctx context.Context, response chatstructs.ChatResponseMessage) {
	sendChatResponse(ctx, response)
}

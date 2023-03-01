package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackTokenRemoveMessage struct {
	TaskID         int                                             `json:"task_id"` //required
	CallbackTokens []MythicRPCCallbackTokenRemoveCallbackTokenData `json:"callbacktokens"`
}
type MythicRPCCallbackTokenRemoveMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type MythicRPCCallbackTokenRemoveCallbackTokenData = agentMessagePostResponseCallbackTokens

func SendMythicRPCCallbackTokenRemove(input MythicRPCCallbackTokenRemoveMessage) (*MythicRPCCallbackTokenRemoveMessageResponse, error) {
	response := MythicRPCCallbackTokenRemoveMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACKTOKEN_REMOVE,
		input,
	); err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	} else if err := json.Unmarshal(responseBytes, &response); err != nil {
		logging.LogError(err, "Failed to parse response back to struct", "response", response)
		return nil, err
	} else {
		return &response, nil
	}
}

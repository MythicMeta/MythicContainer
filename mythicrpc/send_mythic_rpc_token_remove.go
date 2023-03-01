package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCTokenRemoveMessage struct {
	TaskID int                             `json:"task_id"` //required
	Tokens []MythicRPCTokenRemoveTokenData `json:"tokens"`
}
type MythicRPCTokenRemoveMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type MythicRPCTokenRemoveTokenData = agentMessagePostResponseToken

func SendMythicRPCTokenRemove(input MythicRPCTokenRemoveMessage) (*MythicRPCTokenRemoveMessageResponse, error) {
	response := MythicRPCTokenRemoveMessageResponse{}
	for i := 0; i < len(input.Tokens); i++ {
		input.Tokens[i].Action = "remove"
	}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_TOKEN_REMOVE,
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

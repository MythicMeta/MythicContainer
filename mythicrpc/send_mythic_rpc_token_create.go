package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCTokenCreateMessage struct {
	TaskID int                             `json:"task_id"` //required
	Tokens []MythicRPCTokenCreateTokenData `json:"tokens"`
}
type MythicRPCTokenCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type MythicRPCTokenCreateTokenData = agentMessagePostResponseToken

func SendMythicRPCTokenCreate(input MythicRPCTokenCreateMessage) (*MythicRPCTokenCreateMessageResponse, error) {
	response := MythicRPCTokenCreateMessageResponse{}
	for i := 0; i < len(input.Tokens); i++ {
		input.Tokens[i].Action = "add"
	}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_TOKEN_CREATE,
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

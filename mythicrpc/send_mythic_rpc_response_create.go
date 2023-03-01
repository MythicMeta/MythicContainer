package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCResponseCreateMessage struct {
	TaskID   int    `json:"task_id"`
	Response []byte `json:"response"`
}
type MythicRPCResponseCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendMythicRPCResponseCreate(input MythicRPCResponseCreateMessage) (*MythicRPCResponseCreateMessageResponse, error) {
	response := MythicRPCResponseCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_RESPONSE_CREATE,
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

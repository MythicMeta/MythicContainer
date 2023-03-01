package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCPayloadRemoveCommandMessage struct {
	PayloadUUID string   `json:"payload_uuid"` //required
	Commands    []string `json:"commands"`     // required
}
type MythicRPCPayloadRemoveCommandMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendMythicRPCPayloadRemoveCommand(input MythicRPCPayloadRemoveCommandMessage) (*MythicRPCPayloadRemoveCommandMessageResponse, error) {
	response := MythicRPCPayloadRemoveCommandMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_PAYLOAD_REMOVE_COMMAND,
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

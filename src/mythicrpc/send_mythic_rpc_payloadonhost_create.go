package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCPayloadOnHostCreateMessage struct {
	TaskID        int                              `json:"task_id"` //required
	PayloadOnHost MythicRPCPayloadOnHostCreateData `json:"payload_on_host"`
}
type MythicRPCPayloadOnHostCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type MythicRPCPayloadOnHostCreateData struct {
	Host        string  `json:"host"`
	PayloadId   *int    `json:"payload_id"`
	PayloadUUID *string `json:"payload_uuid"`
}

func SendMythicRPCPayloadOnHostCreate(input MythicRPCPayloadOnHostCreateMessage) (*MythicRPCPayloadOnHostCreateMessageResponse, error) {
	response := MythicRPCPayloadOnHostCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_PAYLOADONHOST_CREATE,
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

package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCPayloadCreateFromUUIDMessage struct {
	PayloadUUID    string  `json:"uuid"`
	TaskID         int     `json:"task_id"`
	NewDescription *string `json:"new_description"`
	NewFilename    *string `json:"new_filename"`
	RemoteHost     *string `json:"remote_host"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCPayloadCreateFromUUIDMessageResponse struct {
	Success        bool   `json:"success"`
	Error          string `json:"error"`
	NewPayloadUUID string `json:"new_payload_uuid"`
}

func SendMythicRPCPayloadCreateFromUuid(input MythicRPCPayloadCreateFromUUIDMessage) (*MythicRPCPayloadCreateFromUUIDMessageResponse, error) {
	response := MythicRPCPayloadCreateFromUUIDMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_PAYLOAD_CREATE_FROM_UUID,
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

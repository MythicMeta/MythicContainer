package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackDecryptBytesMessage struct {
	AgentCallbackUUID string `json:"agent_callback_id"` // required if IncludesUUID is false
	Message           []byte `json:"message"`           // required
	IncludesUUID      bool   `json:"include_uuid"`
	IsBase64Encoded   bool   `json:"base64_message"`
}
type MythicRPCCallbackDecryptBytesMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message []byte `json:"message"`
}

func SendMythicRPCCallbackDecryptBytes(input MythicRPCCallbackDecryptBytesMessage) (*MythicRPCCallbackDecryptBytesMessageResponse, error) {
	response := MythicRPCCallbackDecryptBytesMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_DECRYPT_BYTES,
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

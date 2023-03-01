package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackEncryptBytesMessage struct {
	AgentCallbackUUID   string `json:"agent_callback_id"` //required
	Message             []byte `json:"message"`           // required
	IncludeUUID         bool   `json:"include_uuid"`
	Base64ReturnMessage bool   `json:"base64_message"`
}
type MythicRPCCallbackEncryptBytesMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message []byte `json:"message"`
}

func SendMythicRPCCallbackEncryptBytes(input MythicRPCCallbackEncryptBytesMessage) (*MythicRPCCallbackEncryptBytesMessageResponse, error) {
	response := MythicRPCCallbackEncryptBytesMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_ENCRYPT_BYTES,
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

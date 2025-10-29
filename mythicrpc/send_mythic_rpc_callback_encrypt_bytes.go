package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackEncryptBytesMessage struct {
	// AgentCallbackID (Required) - The UUID for the callback that will encrypt the message. Can also be a PayloadUUID or StagingUUID, but requires the C2Profile field as well
	AgentCallbackID string `json:"agent_callback_id"` //required
	// Message (Required) - The actual encrypted message you want to encrypt
	Message []byte `json:"message"`
	// IncludeUUID (Optional) - Should the encrypted message include the UUID in front?
	IncludeUUID bool `json:"include_uuid"`
	// Base64ReturnMessage (Optional) - Should the resulting Message be base64 encoded or left as raw bytes?
	Base64ReturnMessage bool `json:"base64_message"`
	// C2Profile (optional) - If using a Payload UUID then the C2 Profile is required
	C2Profile string `json:"c2_profile"`
}
type MythicRPCCallbackEncryptBytesMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message []byte `json:"message"`
}

// SendMythicRPCCallbackEncryptBytes - Ask Mythic to encrypt a message for a specific callback UUID.
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

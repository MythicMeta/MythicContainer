package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackDecryptBytesMessage struct {
	// AgentCallbackUUID (Required) - The UUID for the callback that will decrypt the message. Can also be a PayloadUUID or StagingUUID, but requires the C2Profile field as well
	AgentCallbackUUID string `json:"agent_callback_id"`
	// Message (Required) - The actual encrypted message you want to decrypt
	Message []byte `json:"message"`
	// IncludesUUID (Optional) - Does the Message include the UUID or not?
	IncludesUUID bool `json:"include_uuid"`
	// IsBase64Encoded (Optional) - Is the Message base64 encoded, or is it just the raw bytes?
	IsBase64Encoded bool `json:"base64_message"`
	// C2Profile (optional) - If using a Payload UUID then the C2 Profile is required
	C2Profile string `json:"c2_profile"`
}
type MythicRPCCallbackDecryptBytesMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message []byte `json:"message"`
}

// SendMythicRPCCallbackDecryptBytes - Ask Mythic to look up the associated callback and decrypt a message for that callback
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

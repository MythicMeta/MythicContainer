package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackCreateMessage struct {
	// PayloadUUID (Required) - What is the UUID of the payload that this new callback will be based on
	PayloadUUID string `json:"payload_uuid"`
	// C2ProfileName (Required) - What is the name of the C2 Profile that this agent is communicating over.
	C2ProfileName string `json:"c2_profile"`
	// EncryptionKey (Optional) - Specify a custom encryption key for use with this callback instead of the
	// C2 profile/Payload's encryption keys.
	EncryptionKey *[]byte `json:"encryption_key"`
	// DecryptionKey (Optional) - Specify a custom decryption key for use with this callback instead of the
	// C2 profile/Payload's decryption keys
	DecryptionKey *[]byte `json:"decryption_key"`
	// CryptoType (Optional) - What kind of crypto is being used? aes256_hmac? none? something else?
	CryptoType string `json:"crypto_type"`
	// User (Optional) - What is the username associated with this new callback
	User string `json:"user"`
	// Host (Optional) - What is the hostname associated with this new callback
	Host string `json:"host"`
	// PID (Optional) - What is the PID associated with this new callback
	PID int `json:"pid"`
	// ExtraInfo (Optional) - Additional information you can store with the callback for context or tracking
	ExtraInfo string `json:"extra_info"`
	// SleepInfo (Optional) - Additional context information about the current sleep data for this callback
	SleepInfo string `json:"sleep_info"`
	// Ip (Optional) - The IP associated with this callback. Use this if you just want to set a single IP address for the callback.
	Ip string `json:"ip"`
	// IPs (Optional) - The array of IP addresses associated with this callback. Use this if you have multiple IP addresses
	// for the callback and want to return them all for the operator to view
	IPs []string `json:"ips"`
	// ExternalIP (Optional) - The external IP address associated with this callback
	ExternalIP string `json:"external_ip"`
	// IntegrityLevel (Optional) - The integrity level associated with the callback.
	// 0 is Unknown, 1 is Low integrity, 2 is Medium integrity, 3 is High integrity, 4 is SYSTEM integrity.
	// 3 and above will result in a red interact button (if you're root, you should return 3+).
	IntegrityLevel int `json:"integrity_level"`
	// Os (Optional) - More detailed OS information than simply the "Windows", "macOS", "Linux", etc associated
	// with the payload
	Os string `json:"os"`
	// Domain (Optional) - the domain associated with the callback.
	Domain string `json:"domain"`
	// Architecture (Optional) - The architecture of the callback (x86, x64, arm64, etc)
	Architecture string `json:"architecture"`
	// Description (Optional) - Set a description for the new callback
	Description string `json:"description"`
	// ProcessName (Optional) - The name of process associated with the new callback.
	ProcessName string `json:"process_name" mapstructure:"process_name"`
}
type MythicRPCCallbackCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	// CallbackUUID - The AgentCallbackID for the new callback that was created.
	CallbackUUID string `json:"callback_uuid"`
}

// SendMythicRPCCallbackCreate - Register a new callback within Mythic
func SendMythicRPCCallbackCreate(input MythicRPCCallbackCreateMessage) (*MythicRPCCallbackCreateMessageResponse, error) {
	response := MythicRPCCallbackCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_CREATE,
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

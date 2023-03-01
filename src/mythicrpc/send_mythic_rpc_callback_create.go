package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackCreateMessage struct {
	PayloadUUID    string  `json:"payload_uuid"` // required
	C2ProfileName  string  `json:"c2_profile"`   // required
	EncryptionKey  *[]byte `json:"encryption_key"`
	DecryptionKey  *[]byte `json:"decryption_key"`
	CryptoType     string  `json:"crypto_type"`
	User           string  `json:"user"`
	Host           string  `json:"host"`
	PID            int     `json:"pid"`
	ExtraInfo      string  `json:"extra_info"`
	SleepInfo      string  `json:"sleep_info"`
	Ip             string  `json:"ip"`
	ExternalIP     string  `json:"external_ip"`
	IntegrityLevel int     `json:"integrity_level"`
	Os             string  `json:"os"`
	Domain         string  `json:"domain"`
	Architecture   string  `json:"architecture"`
}
type MythicRPCCallbackCreateMessageResponse struct {
	Success      bool   `json:"success"`
	Error        string `json:"error"`
	CallbackUUID string `json:"callback_uuid"`
}

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

package mythicrpc

import (
	"encoding/json"
	"time"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackNextCheckinRangeMessage struct {
	SleepInterval int       `json:"sleep_interval"`
	SleepJitter   int       `json:"sleep_jitter"`
	LastCheckin   time.Time `json:"last_checkin"`
}

type MythicRPCCallbackNextCheckinRangeMessageResponse struct {
	Success bool      `json:"success"`
	Error   string    `json:"error"`
	Min     time.Time `json:"min"`
	Max     time.Time `json:"max"`
}

// SendMythicRPCCallbackEncryptBytes - Ask Mythic to encrypt a message for a specific callback UUID.
func SendMythicRPCCallbackNextCheckinRange(input MythicRPCCallbackNextCheckinRangeMessage) (*MythicRPCCallbackNextCheckinRangeMessageResponse, error) {
	response := MythicRPCCallbackNextCheckinRangeMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_NEXT_CHECKIN_RANGE,
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

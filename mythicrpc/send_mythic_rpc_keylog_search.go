package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCKeylogSearchMessage struct {
	TaskID        int                             `json:"task_id"` //required
	SearchKeylogs MythicRPCKeylogSearchKeylogData `json:"keylogs"`
}
type MythicRPCKeylogSearchMessageResponse struct {
	Success bool                              `json:"success"`
	Error   string                            `json:"error"`
	Keylogs []MythicRPCKeylogSearchKeylogData `json:"keylogs"`
}
type MythicRPCKeylogSearchKeylogData struct {
	User        *string `json:"user" `         // optional
	WindowTitle *string `json:"window_title" ` // optional
	Keystrokes  *[]byte `json:"keystrokes" `   // optional
}

func SendMythicRPCKeylogSearch(input MythicRPCKeylogSearchMessage) (*MythicRPCKeylogSearchMessageResponse, error) {
	response := MythicRPCKeylogSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_KEYLOG_SEARCH,
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

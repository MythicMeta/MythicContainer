package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCC2UpdateStatusMessage struct {
	C2Profile             string `json:"c2_profile"`     // required
	InternalServerRunning bool   `json:"server_running"` // required
	Error                 string `json:"error"`
}
type MythicRPCC2UpdateStatusMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// SendMythicRPCCallbackCreate - Register a new callback within Mythic
func SendMythicRPCC2UpdateStatus(input MythicRPCC2UpdateStatusMessage) (*MythicRPCC2UpdateStatusMessageResponse, error) {
	response := MythicRPCC2UpdateStatusMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_C2_UPDATE_STATUS,
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

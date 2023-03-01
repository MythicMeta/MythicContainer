package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCResponseSearchMessage struct {
	TaskID   int    `json:"task_id"`
	Response string `json:"response"`
}
type MythicRPCResponseSearchMessageResponse struct {
	Success   bool                `json:"success"`
	Error     string              `json:"error"`
	Responses []MythicRPCResponse `json:"responses"`
}
type MythicRPCResponse struct {
	ResponseID int    `json:"response_id"`
	Response   []byte `json:"response"`
	TaskID     int    `json:"task_id"`
}

func SendMythicRPCResponseSearch(input MythicRPCResponseSearchMessage) (*MythicRPCResponseSearchMessageResponse, error) {
	response := MythicRPCResponseSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_RESPONSE_SEARCH,
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

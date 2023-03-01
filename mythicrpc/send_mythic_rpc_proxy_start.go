package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCProxyStartMessage struct {
	TaskID   int    `json:"task_id"`
	Port     int    `json:"port"`
	PortType string `json:"port_type"`
}
type MythicRPCProxyStartMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendMythicRPCProxyStart(input MythicRPCProxyStartMessage) (*MythicRPCProxyStartMessageResponse, error) {
	response := MythicRPCProxyStartMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_PROXY_START,
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

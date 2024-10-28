package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCProxyStopMessage struct {
	TaskID   int    `json:"task_id"`
	Port     int    `json:"port"`
	PortType string `json:"port_type"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type MythicRPCProxyStopMessageResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	LocalPort int    `json:"local_port"`
}

func SendMythicRPCProxyStop(input MythicRPCProxyStopMessage) (*MythicRPCProxyStopMessageResponse, error) {
	response := MythicRPCProxyStopMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_PROXY_STOP,
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

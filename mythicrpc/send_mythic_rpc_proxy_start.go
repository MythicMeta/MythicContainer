package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCProxyStartMessage struct {
	// TaskID - the TaskID that's starting the proxy connection
	TaskID int `json:"task_id"`
	// LocalPort - for SOCKS, this is the port to open on the Mythic server.
	// For interactive, this is the port to open on the Mythic server
	// For rpfwd, this is the port to open on the host where your agent is running.
	LocalPort int `json:"local_port"`
	// RemotePort - This only needs to be set for rpfwd - this is the remote port to connect to when the LocalPort gets a connection
	RemotePort int `json:"remote_port"`
	// RemoteIP - This only needs to be set for rpfwd - this is the remote ip to connect to when the LocalPort gets a connection
	RemoteIP string `json:"remote_ip"`
	// PortType - What type of proxy connection are you opening
	// CALLBACK_PORT_TYPE_SOCKS
	// CALLBACK_PORT_TYPE_RPORTFWD
	// CALLBACK_PORT_TYPE_INTERACTIVE
	PortType string `json:"port_type"`
	// Username - This only can be set for socks - this allows auth for connecting to the port opened on the Mythic server
	Username string `json:"username"`
	// Password - This only can be set for socks - this allows auth for connecting to the port opened on the Mythic server
	Password string `json:"password"`
}
type MythicRPCProxyStartMessageResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	LocalPort int    `json:"local_port"`
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

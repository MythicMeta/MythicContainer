package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MESSAGE_LEVEL = string

const (
	MESSAGE_LEVEL_INFO          MESSAGE_LEVEL = "info"
	MESSAGE_LEVEL_DEBUG         MESSAGE_LEVEL = "debug"
	MESSAGE_LEVEL_AUTH          MESSAGE_LEVEL = "auth"
	MESSAGE_LEVEL_AGENT_MESSAGE MESSAGE_LEVEL = "agent"
	MESSAGE_LEVEL_API           MESSAGE_LEVEL = "api"
)

type MythicRPCOperationEventLogCreateMessage struct {
	// three optional ways to specify the operation
	TaskID          *int    `json:"task_id"`
	CallbackID      *int    `json:"callback_id"`
	AgentCallbackID *string `json:"callback_agent_id"`
	OperationID     *int    `json:"operation_id"`
	// the data to store
	Message      string        `json:"message"`
	Warning      bool          `json:"warning"`
	MessageLevel MESSAGE_LEVEL `json:"level"` //info or warning
}
type MythicRPCOperationEventLogCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendMythicRPCOperationEventLogCreate(input MythicRPCOperationEventLogCreateMessage) (*MythicRPCOperationEventLogCreateMessageResponse, error) {
	response := MythicRPCOperationEventLogCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_EVENTLOG_CREATE,
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

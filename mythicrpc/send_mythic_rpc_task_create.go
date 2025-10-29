package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCTaskCreateMessage struct {
	AgentCallbackID    *string `json:"agent_callback_id"`
	CallbackID         *int    `json:"callback_id"`
	CommandName        string  `json:"command_name"`
	Params             string  `json:"params"`
	ParameterGroupName *string `json:"parameter_group_name,omitempty"`
	Token              *int    `json:"token,omitempty"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCTaskCreateMessageResponse struct {
	Success       bool   `json:"success"`
	Error         string `json:"error"`
	TaskID        int    `json:"task_id"`
	TaskDisplayID int    `json:"task_display_id"`
}

func SendMythicRPCTaskCreate(input MythicRPCTaskCreateMessage) (*MythicRPCTaskCreateMessageResponse, error) {
	response := MythicRPCTaskCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_TASK_CREATE,
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

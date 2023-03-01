package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCTaskCreateSubtaskGroupMessage struct {
	TaskID                int                                    `json:"task_id"`    // required
	GroupName             string                                 `json:"group_name"` // required
	GroupCallbackFunction *string                                `json:"group_callback_function,omitempty"`
	Tasks                 []MythicRPCTaskCreateSubtaskGroupTasks `json:"tasks"` // required

}

type MythicRPCTaskCreateSubtaskGroupTasks struct {
	SubtaskCallbackFunction *string `json:"subtask_callback_function,omitempty"`
	CommandName             string  `json:"command_name"` // required
	Params                  string  `json:"params"`       // required
	ParameterGroupName      *string `json:"parameter_group_name,omitempty"`
	Token                   *int    `json:"token,omitempty"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCTaskCreateSubtaskGroupMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	TaskIDs []int  `json:"task_ids"`
}

func SendMythicRPCTaskCreateSubtaskGroup(input MythicRPCTaskCreateSubtaskGroupMessage) (*MythicRPCTaskCreateSubtaskGroupMessageResponse, error) {
	response := MythicRPCTaskCreateSubtaskGroupMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_TASK_CREATE_SUBTASK_GROUP,
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

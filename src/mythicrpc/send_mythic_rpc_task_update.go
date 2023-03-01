package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCTaskUpdateMessage struct {
	TaskID            int     `json:"task_id"`
	UpdateStatus      *string `json:"update_status,omitempty"`
	UpdateStdout      *string `json:"update_stdout,omitempty"`
	UpdateStderr      *string `json:"update_stderr,omitempty"`
	UpdateCommandName *string `json:"update_command_name,omitempty"`
	UpdateCompleted   *bool   `json:"update_completed,omitempty"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCTaskUpdateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendMythicRPCTaskUpdate(input MythicRPCTaskUpdateMessage) (*MythicRPCTaskUpdateMessageResponse, error) {
	response := MythicRPCTaskUpdateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_TASK_UPDATE,
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

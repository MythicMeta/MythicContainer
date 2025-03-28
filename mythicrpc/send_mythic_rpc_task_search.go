package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCTaskSearchMessage struct {
	TaskID              int       `json:"task_id"`
	SearchTaskID        *int      `json:"search_task_id"`
	SearchTaskDisplayID *int      `json:"search_task_display_id"`
	SearchAgentTaskID   *string   `json:"agent_task_id,omitempty"`
	SearchHost          *string   `json:"host,omitempty"`
	SearchCallbackID    *int      `json:"callback_id,omitempty"`
	SearchCompleted     *bool     `json:"completed,omitempty"`
	SearchCommandNames  *[]string `json:"command_names,omitempty"`
	SearchParams        *string   `json:"params,omitempty"`
	SearchParentTaskID  *int      `json:"parent_task_id,omitempty"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCTaskSearchMessageResponse struct {
	Success bool                    `json:"success"`
	Error   string                  `json:"error"`
	Tasks   []PTTaskMessageTaskData `json:"tasks"`
}

func SendMythicRPCTaskSearch(input MythicRPCTaskSearchMessage) (*MythicRPCTaskSearchMessageResponse, error) {
	response := MythicRPCTaskSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_TASK_SEARCH,
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

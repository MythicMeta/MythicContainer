package mythicrpc

import (
	"context"
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCDirectFileTokenCreateMessage struct {
	AgentTaskID         *string `json:"agent_task_id"`
	AgentCallbackID     *string `json:"agent_callback_id"`
	PayloadUUID         *string `json:"payload_uuid"`
	AgentFileID         string  `json:"agent_file_id"`
	APITokenID          *int    `json:"apitoken_id"`
	EventStepInstanceID *int    `json:"eventstep_instance_id"`
	Action              string  `json:"action"`
}

type MythicRPCDirectFileTokenCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Token   string `json:"token"`
}

func SendMythicRPCDirectFileTokenCreate(ctx context.Context, input MythicRPCDirectFileTokenCreateMessage) (*MythicRPCDirectFileTokenCreateMessageResponse, error) {
	response := MythicRPCDirectFileTokenCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessageWithContext(
		ctx,
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_DIRECT_FILE_TOKEN_CREATE,
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

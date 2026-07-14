package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/MythicMeta/MythicContainer/logging"
)

type DirectFileTokenCreateMessage struct {
	AgentTaskID         *string `json:"agent_task_id"`
	AgentCallbackID     *string `json:"agent_callback_id"`
	PayloadUUID         *string `json:"payload_uuid"`
	AgentFileID         string  `json:"agent_file_id"`
	APITokenID          *int    `json:"apitoken_id"`
	EventStepInstanceID *int    `json:"eventstep_instance_id"`
	Action              string  `json:"action"`
}

type DirectFileTokenCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Token   string `json:"token"`
}

func SendDirectFileTokenCreate(ctx context.Context, input DirectFileTokenCreateMessage) (*DirectFileTokenCreateMessageResponse, error) {
	response := DirectFileTokenCreateMessageResponse{}
	responseBytes, err := RabbitMQConnection.SendRPCStructMessageWithContext(
		ctx,
		MYTHIC_EXCHANGE,
		MYTHIC_RPC_DIRECT_FILE_TOKEN_CREATE,
		input,
	)
	if err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	}
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		logging.LogError(err, "Failed to parse response back to struct", "response", response)
		return nil, err
	}
	return &response, nil
}

func RequestDirectFileToken(ctx context.Context, fileID string, action string) (string, error) {
	response, err := SendDirectFileTokenCreate(ctx, DirectFileTokenCreateMessage{
		AgentFileID: fileID,
		Action:      action,
	})
	if err != nil {
		return "", err
	}
	if !response.Success {
		if response.Error != "" {
			return "", errors.New(response.Error)
		}
		return "", errors.New("failed to create direct file token")
	}
	if response.Token == "" {
		return "", errors.New("direct file token response did not include a token")
	}
	return response.Token, nil
}

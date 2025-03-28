package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCTagCreateMessage struct {
	TagTypeID    int                    `json:"tagtype_id"`
	URL          string                 `json:"url"`
	Source       string                 `json:"source"`
	Data         map[string]interface{} `json:"data"`
	TaskID       *int                   `json:"task_id"`
	FileID       *int                   `json:"file_id"`
	CredentialID *int                   `json:"credential_id"`
	MythicTreeID *int                   `json:"mythic_tree_id"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCTagCreateMessageResponse struct {
	Success bool             `json:"success"`
	Error   string           `json:"error"`
	Tag     MythicRPCTagData `json:"tag"`
}

func SendMythicRPCTagCreate(input MythicRPCTagCreateMessage) (*MythicRPCTagCreateMessageResponse, error) {
	response := MythicRPCTagCreateMessageResponse{}
	responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_TAG_CREATE,
		input,
	)
	if err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		logging.LogError(err, "Failed to parse response back to struct", "response", response)
		return nil, err
	}
	return &response, nil
}

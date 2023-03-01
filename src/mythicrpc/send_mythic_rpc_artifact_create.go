package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCArtifactCreateMessage struct {
	TaskID           int     `json:"task_id"`
	ArtifactMessage  string  `json:"message"`
	BaseArtifactType string  `json:"base_artifact"`
	ArtifactHost     *string `json:"host,omitempty"`
}
type MythicRPCArtifactCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendMythicRPCArtifactCreate(input MythicRPCArtifactCreateMessage) (*MythicRPCArtifactCreateMessageResponse, error) {
	response := MythicRPCArtifactCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_ARTIFACT_CREATE,
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

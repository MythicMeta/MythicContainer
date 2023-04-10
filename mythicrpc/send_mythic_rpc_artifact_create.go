package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCArtifactCreateMessage struct {
	// TaskID (Required) - the task associated with this new artifact for Mythic to track
	TaskID int `json:"task_id"`
	// ArtifactMessage (Required) - the actual artifact string you want to store
	ArtifactMessage string `json:"message"`
	// BaseArtifactType (Required) - what kind of artifact is it? Process Create? File Removal? etc
	BaseArtifactType string `json:"base_artifact"`
	// ArtifactHost (Optional) - what's the hostname for where this artifact happened? If none is specified, it's assumed to be the same host where the task ran
	ArtifactHost *string `json:"host,omitempty"`
}
type MythicRPCArtifactCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// SendMythicRPCArtifactCreate - Create a new artifact for Mythic to track.
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

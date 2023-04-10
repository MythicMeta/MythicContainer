package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCArtifactSearchMessage struct {
	// TaskID (Required) - What is the current task that's searching for artifact information.
	TaskID int `json:"task_id"`
	// SearchArtifacts (Required) - Additional structure of data used to search artifacts.
	SearchArtifacts MythicRPCArtifactearchArtifactData `json:"artifact"`
}
type MythicRPCArtifactSearchMessageResponse struct {
	Success   bool                                 `json:"success"`
	Error     string                               `json:"error"`
	Artifacts []MythicRPCArtifactearchArtifactData `json:"artifacts"`
}
type MythicRPCArtifactearchArtifactData struct {
	// Host (Optional) - When searching, you can filter your artifacts by the hostname.
	// As a response, this will always be populated.
	Host *string `json:"host" ` // optional
	// ArtifactType (Optional) - When searching, you can filter your artifacts by the type of artifact.
	// As a response, this will always be populated.
	ArtifactType *string `json:"artifact_type"` //optional
	// ArtifactMessage (Optional) - When searching, you can filter your artifacts by what the message contains.
	// As a response, this will always be populated.
	ArtifactMessage *string `json:"artifact_message"` //optional
	// TaskID (Optional) - When searching, you can filter your artifacts to those created by a certain task.
	// As a response, this will always be populated.
	TaskID *int `json:"task_id"` //optional
}

// SendMythicRPCArtifactSearch - Search for artifacts that are tracked by Mythic.
func SendMythicRPCArtifactSearch(input MythicRPCArtifactSearchMessage) (*MythicRPCArtifactSearchMessageResponse, error) {
	response := MythicRPCArtifactSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_ARTIFACT_SEARCH,
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

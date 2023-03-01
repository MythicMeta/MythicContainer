package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCArtifactSearchMessage struct {
	TaskID          int                                `json:"task_id"` //required
	SearchArtifacts MythicRPCArtifactearchArtifactData `json:"artifact"`
}
type MythicRPCArtifactSearchMessageResponse struct {
	Success   bool                                 `json:"success"`
	Error     string                               `json:"error"`
	Artifacts []MythicRPCArtifactearchArtifactData `json:"artifacts"`
}
type MythicRPCArtifactearchArtifactData struct {
	Host            *string `json:"host" `            // optional
	ArtifactType    *string `json:"artifact_type"`    //optional
	ArtifactMessage *string `json:"artifact_message"` //optional
	TaskID          *int    `json:"task_id"`          //optional
}

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

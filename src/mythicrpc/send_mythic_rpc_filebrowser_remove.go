package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCFileBrowserRemoveMessage struct {
	TaskID       int                                         `json:"task_id"` //required
	RemovedFiles []MythicRPCFileBrowserRemoveFileBrowserData `json:"removed_files"`
}
type MythicRPCFileBrowserRemoveMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type MythicRPCFileBrowserRemoveFileBrowserData = agentMessagePostResponseRemovedFiles
type agentMessagePostResponseRemovedFiles struct {
	Host *string `json:"host,omitempty" mapstructure:"host,omitempty"`
	Path string  `json:"path" mapstructure:"path"` // full path to file removed
}

func SendMythicRPCFileBrowserRemove(input MythicRPCFileBrowserRemoveMessage) (*MythicRPCFileBrowserRemoveMessageResponse, error) {
	response := MythicRPCFileBrowserRemoveMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_FILEBROWSER_REMOVE,
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

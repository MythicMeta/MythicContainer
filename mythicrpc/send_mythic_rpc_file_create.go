package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
	"github.com/MythicMeta/MythicContainer/utils/mythicutils"
)

// MythicRPCFileCreateMessage Must supply one of TaskID, PayloadUUID, or AgentCallbackID so that Mythic can track the file for the right operation
type MythicRPCFileCreateMessage struct {
	TaskID              int    `json:"task_id"`
	PayloadUUID         string `json:"payload_uuid"`
	AgentCallbackID     string `json:"agent_callback_id"`
	FileContents        []byte `json:"-"`
	DeleteAfterFetch    bool   `json:"delete_after_fetch"`
	Filename            string `json:"filename"`
	IsScreenshot        bool   `json:"is_screenshot"`
	IsDownloadFromAgent bool   `json:"is_download"`
	RemotePathOnTarget  string `json:"remote_path"`
	TargetHostName      string `json:"host"`
	Comment             string `json:"comment"`
}
type MythicRPCFileCreateMessageResponse struct {
	Success     bool   `json:"success"`
	Error       string `json:"error"`
	AgentFileID string `json:"agent_file_id"`
}

func SendMythicRPCFileCreate(input MythicRPCFileCreateMessage) (*MythicRPCFileCreateMessageResponse, error) {
	response := MythicRPCFileCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_FILE_CREATE,
		input,
	); err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	} else if err := json.Unmarshal(responseBytes, &response); err != nil {
		logging.LogError(err, "Failed to parse response back to struct", "response", response)
		return nil, err
	} else if response.Success {
		if err := mythicutils.SendFileToMythic(&input.FileContents, response.AgentFileID); err != nil {
			logging.LogError(err, "Failed to send file contents to Mythic")
			return nil, err
		} else {
			return &response, nil
		}
	} else {
		return &response, nil
	}
}

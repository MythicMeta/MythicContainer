package mythicrpc

import (
	"encoding/json"
	"time"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCFileSearchMessage struct {
	TaskID              int    `json:"task_id"`
	CallbackID          int    `json:"callback_id"`
	Filename            string `json:"filename"`
	LimitByCallback     bool   `json:"limit_by_callback"`
	MaxResults          int    `json:"max_results"`
	Comment             string `json:"comment"`
	AgentFileID         string `json:"file_id"`
	IsPayload           bool   `json:"is_payload"`
	IsDownloadFromAgent bool   `json:"is_download_from_agent"`
	IsScreenshot        bool   `json:"is_screenshot"`
}
type FileData struct {
	AgentFileId         string    `json:"agent_file_id"`
	Filename            string    `json:"filename"`
	Comment             string    `json:"comment"`
	Complete            bool      `json:"complete"`
	IsPayload           bool      `json:"is_payload"`
	IsDownloadFromAgent bool      `json:"is_download_from_agent"`
	IsScreenshot        bool      `json:"is_screenshot"`
	FullRemotePath      string    `json:"full_remote_path"`
	Host                string    `json:"host"`
	TaskID              int       `json:"task_id"`
	Md5                 string    `json:"md5"`
	Sha1                string    `json:"sha1"`
	Timestamp           time.Time `json:"timestamp"`
	Command             string    `json:"cmd"`
	Tags                []string  `json:"tags"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCFileSearchMessageResponse struct {
	Success bool       `json:"success"`
	Error   string     `json:"error"`
	Files   []FileData `json:"files"`
}

func SendMythicRPCFileSearch(input MythicRPCFileSearchMessage) (*MythicRPCFileSearchMessageResponse, error) {
	response := MythicRPCFileSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_FILE_SEARCH,
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

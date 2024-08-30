package loggingstructs

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"time"
)

type NewFileLog struct {
	loggingMessageBase
	Data NewFileData `json:"data"`
}
type NewFileData struct {
	ID                  int       `db:"id" json:"id" mapstructure:"id"`
	AgentFileID         string    `db:"agent_file_id" json:"agent_file_id" mapstructure:"agent_file_id"`
	TotalChunks         int       `db:"total_chunks" json:"total_chunks" mapstructure:"total_chunks"`
	ChunksReceived      int       `db:"chunks_received" json:"chunks_received" mapstructure:"chunks_received"`
	ChunkSize           int       `db:"chunk_size" json:"chunk_size" mapstructure:"chunk_size"`
	TaskID              *int      `db:"task_id" json:"task_id" mapstructure:"task_id"`
	Complete            bool      `db:"complete" json:"complete" mapstructure:"complete"`
	Path                string    `db:"path" json:"path" mapstructure:"path"`
	FullRemotePath      []byte    `db:"full_remote_path" json:"full_remote_path" mapstructure:"full_remote_path"`
	Host                string    `db:"host" json:"host" mapstructure:"host"`
	IsPayload           bool      `db:"is_payload" json:"is_payload" mapstructure:"is_payload"`
	IsScreenshot        bool      `db:"is_screenshot" json:"is_screenshot" mapstructure:"is_screenshot"`
	IsDownloadFromAgent bool      `db:"is_download_from_agent" json:"is_download_from_agent" mapstructure:"is_download_from_agent"`
	MythicTreeID        *int      `db:"mythictree_id" json:"mythictree_id" mapstructure:"mythictree_id"`
	Filename            []byte    `db:"filename" json:"filename" mapstructure:"filename"`
	DeleteAfterFetch    bool      `db:"delete_after_fetch" json:"delete_after_fetch" mapstructure:"delete_after_fetch"`
	OperationID         int       `db:"operation_id" json:"operation_id" mapstructure:"operation_id"`
	Timestamp           time.Time `db:"timestamp" json:"timestamp" mapstructure:"timestamp"`
	Deleted             bool      `db:"deleted" json:"deleted" mapstructure:"deleted"`
	OperatorID          int       `db:"operator_id" json:"operator_id" mapstructure:"operator_id"`
	Md5                 string    `db:"md5" json:"md5" mapstructure:"md5"`
	Sha1                string    `db:"sha1" json:"sha1" mapstructure:"sha1"`
	Comment             string    `db:"comment" json:"comment" mapstructure:"comment"`
}

func init() {
	AllLoggingData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         LOG_TYPE_FILE,
		RabbitmqProcessingFunction: processNewFileLog,
	})
}

func processNewFileLog(input []byte) {
	inputStruct := NewFileLog{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process message")
	} else {
		for _, webhook := range AllLoggingData.GetAllNames() {
			if AllLoggingData.Get(webhook).GetLoggingDefinition().NewFileFunction != nil {
				AllLoggingData.Get(webhook).GetLoggingDefinition().NewFileFunction(inputStruct)
			}
		}
	}
}

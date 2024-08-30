package loggingstructs

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"time"
)

type NewArtifactLog struct {
	loggingMessageBase
	Data NewArtifactLogData `json:"data"`
}
type NewArtifactLogData struct {
	ID           int       `json:"id"`
	TaskID       *int      `json:"task_id,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	Artifact     []byte    `json:"artifact"`
	BaseArtifact string    `json:"base_artifact"`
	OperationID  int       `json:"operation_id"`
	Host         string    `json:"host"`
}

func init() {
	AllLoggingData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         LOG_TYPE_ARTIFACT,
		RabbitmqProcessingFunction: processNewArtifactLog,
	})
}

func processNewArtifactLog(input []byte) {
	inputStruct := NewArtifactLog{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process message")
	} else {
		for _, webhook := range AllLoggingData.GetAllNames() {
			if AllLoggingData.Get(webhook).GetLoggingDefinition().NewArtifactFunction != nil {
				AllLoggingData.Get(webhook).GetLoggingDefinition().NewArtifactFunction(inputStruct)
			}
		}
	}
}

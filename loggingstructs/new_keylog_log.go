package loggingstructs

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"time"
)

type NewKeylogLog struct {
	loggingMessageBase
	Data NewKeylogLogData `json:"data"`
}
type NewKeylogLogData struct {
	ID          int       `json:"id" mapstructure:"id"`
	TaskID      int       `json:"task_id" mapstructure:"task_id"`
	Keystrokes  []byte    `json:"keystrokes" mapstructure:"keystrokes"`
	Window      string    `json:"window" mapstructure:"window"`
	Timestamp   time.Time `json:"timestamp" mapstructure:"timestamp"`
	OperationID int       `json:"operation_id" mapstructure:"operation_id"`
	User        string    `json:"user" mapstructure:"user"`
}

func init() {
	AllLoggingData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         LOG_TYPE_KEYLOG,
		RabbitmqProcessingFunction: processNewKeylogLog,
	})
}

func processNewKeylogLog(input []byte) {
	inputStruct := NewKeylogLog{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process message")
	} else {
		for _, webhook := range AllLoggingData.GetAllNames() {
			if AllLoggingData.Get(webhook).GetLoggingDefinition().NewKeylogFunction != nil {
				AllLoggingData.Get(webhook).GetLoggingDefinition().NewKeylogFunction(inputStruct)
			}
		}
	}
}

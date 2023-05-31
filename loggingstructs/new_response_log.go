package loggingstructs

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/logging"
)

type NewResponseLog struct {
	loggingMessageBase
	Data ResponseLogData `json:"data"`
}

type ResponseLogData struct {
	ID            int    `json:"id" mapstructure:"id"`
	Response      []byte `json:"response" mapstructure:"response"`
	TaskID        int    `json:"task_id" mapstructure:"task_id"`
	TaskDisplayID int    `json:"task_display_id" mapstructure:"task_display_id"`
	Timestamp     string `json:"timestamp" mapstructure:"timestamp"`
}

func init() {
	AllLoggingData.Get("").AddDirectMethod(RabbitmqDirectMethod{
		RabbitmqRoutingKey:         LOG_TYPE_RESPONSE,
		RabbitmqProcessingFunction: processResponseLog,
	})
}

func processResponseLog(input []byte) {
	inputStruct := NewResponseLog{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process message")
	} else {
		for _, webhook := range AllLoggingData.GetAllNames() {
			if AllLoggingData.Get(webhook).GetLoggingDefinition().NewResponseFunction != nil {
				AllLoggingData.Get(webhook).GetLoggingDefinition().NewResponseFunction(inputStruct)
			}
		}
	}
}

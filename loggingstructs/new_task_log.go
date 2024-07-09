package loggingstructs

import (
	"encoding/json"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

type NewTaskLog struct {
	loggingMessageBase
	Data NewTaskLogData `json:"data"`
}
type NewTaskLogData = agentstructs.PTTaskMessageTaskData

func init() {
	AllLoggingData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         LOG_TYPE_TASK,
		RabbitmqProcessingFunction: processNewTaskLog,
	})
}

func processNewTaskLog(input []byte) {
	inputStruct := NewTaskLog{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process message")
	} else {
		for _, webhook := range AllLoggingData.GetAllNames() {
			if AllLoggingData.Get(webhook).GetLoggingDefinition().NewTaskFunction != nil {
				AllLoggingData.Get(webhook).GetLoggingDefinition().NewTaskFunction(inputStruct)
			}
		}
	}
}

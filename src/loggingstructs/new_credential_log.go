package loggingstructs

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/logging"
	"time"
)

type NewCredentialLog struct {
	loggingMessageBase
	Data NewCredentialLogData `json:"data"`
}
type NewCredentialLogData struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"`
	TaskID      *int      `json:"task_id"`
	Account     string    `json:"account"`
	Realm       string    `json:"realm"`
	OperationID int       `json:"operation_id"`
	Timestamp   time.Time `json:"timestamp"`
	Credential  string    `json:"credential"`
	OperatorID  int       `json:"operator_id"`
	Comment     string    `json:"comment"`
	Deleted     bool      `json:"deleted"`
	Metadata    string    `json:"metadata"`
}

func init() {
	AllLoggingData.Get("").AddDirectMethod(RabbitmqDirectMethod{
		RabbitmqRoutingKey:         LOG_TYPE_CREDENTIAL,
		RabbitmqProcessingFunction: processNewCredentialLog,
	})
}

func processNewCredentialLog(input []byte) {
	inputStruct := NewCredentialLog{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process message")
	} else {
		for _, webhook := range AllLoggingData.GetAllNames() {
			if AllLoggingData.Get(webhook).GetLoggingDefinition().NewCredentialFunction != nil {
				AllLoggingData.Get(webhook).GetLoggingDefinition().NewCredentialFunction(inputStruct)
			}
		}
	}
}

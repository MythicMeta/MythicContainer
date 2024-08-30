package webhookstructs

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

type NewFeedbackWebookMessage struct {
	webhookMessageBase
	Data NewFeedbackWebhookData `json:"data"`
}

type NewFeedbackWebhookData struct {
	TaskID        *int   `json:"task_id,omitempty"`
	TaskDisplayID *int   `json:"display_id,omitempty"`
	Message       string `json:"message"`
	FeedbackType  string `json:"feedback_type"`
}

// Register this method with rabbitmq so it can be called
func init() {
	AllWebhookData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         WEBHOOK_TYPE_NEW_FEEDBACK,
		RabbitmqProcessingFunction: processNewFeedbackWebhook,
	})
}

func processNewFeedbackWebhook(input []byte) {
	inputStruct := NewFeedbackWebookMessage{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process new feedback webhook message")
	} else {
		for _, webhook := range AllWebhookData.GetAllNames() {
			if AllWebhookData.Get(webhook).GetWebhookDefinition().NewFeedbackFunction != nil {
				AllWebhookData.Get(webhook).GetWebhookDefinition().NewFeedbackFunction(inputStruct)
			}
		}
	}
}

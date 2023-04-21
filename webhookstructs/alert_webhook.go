package webhookstructs

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/logging"
)

type NewAlertWebhookMessage struct {
	webhookMessageBase
	Data NewAlertWebhookData `json:"data"`
}
type NewAlertWebhookData struct {
	OperatorID int    `json:"operator_id"`
	Message    string `json:"message"`
	Source     string `json:"source"`
	Count      int    `json:"count"`
	Timestamp  string `json:"timestamp"`
}

// Register this method with rabbitmq so it can be called
func init() {
	AllWebhookData.Get("").AddDirectMethod(RabbitmqDirectMethod{
		RabbitmqRoutingKey:         WEBHOOK_TYPE_NEW_ALERT,
		RabbitmqProcessingFunction: processNewAlertWebhook,
	})
}
func processNewAlertWebhook(input []byte) {
	inputStruct := NewAlertWebhookMessage{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process new callback webhook message")
	} else {
		// success, so do RPC calls to Mythic to get more context or send off webhook now
		for _, webhook := range AllWebhookData.GetAllNames() {
			if AllWebhookData.Get(webhook).GetWebhookDefinition().NewAlertFunction != nil {
				AllWebhookData.Get(webhook).GetWebhookDefinition().NewAlertFunction(inputStruct)
			}
		}
	}
}

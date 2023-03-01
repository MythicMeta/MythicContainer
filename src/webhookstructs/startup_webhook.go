package webhookstructs

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/logging"
)

type NewStartupWebhookMessage struct {
	webhookMessageBase
	Data NewStartupWebhookData `json:"data"`
}
type NewStartupWebhookData struct {
	StartupMessage string `json:"startup_message"`
}

// Register this method with rabbitmq so it can be called
func init() {
	AllWebhookData.Get("").AddDirectMethod(RabbitmqDirectMethod{
		RabbitmqRoutingKey:         WEBHOOK_TYPE_NEW_STARTUP,
		RabbitmqProcessingFunction: processNewStartupWebhook,
	})
}
func processNewStartupWebhook(input []byte) {
	inputStruct := NewStartupWebhookMessage{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process new callback webhook message")
	} else {
		// success, so do RPC calls to Mythic to get more context or send off webhook now
		for _, webhook := range AllWebhookData.GetAllNames() {
			if AllWebhookData.Get(webhook).GetWebhookDefinition().NewStartupFunction != nil {
				AllWebhookData.Get(webhook).GetWebhookDefinition().NewStartupFunction(inputStruct)
			}
		}
	}
}

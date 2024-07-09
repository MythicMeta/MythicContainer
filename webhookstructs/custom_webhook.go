package webhookstructs

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

type NewCustomWebhookMessage struct {
	webhookMessageBase
	Data map[string]string `json:"data"`
}

// Register this method with rabbitmq so it can be called
func init() {
	AllWebhookData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         WEBHOOK_TYPE_NEW_CUSTOM,
		RabbitmqProcessingFunction: processNewCustomWebhook,
	})
}
func processNewCustomWebhook(input []byte) {
	inputStruct := NewCustomWebhookMessage{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process new callback webhook message")
	} else {
		// success, so do RPC calls to Mythic to get more context or send off webhook now
		for _, webhook := range AllWebhookData.GetAllNames() {
			if AllWebhookData.Get(webhook).GetWebhookDefinition().NewCustomFunction != nil {
				AllWebhookData.Get(webhook).GetWebhookDefinition().NewCustomFunction(inputStruct)
			}
		}
	}
}

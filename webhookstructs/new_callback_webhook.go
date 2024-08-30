package webhookstructs

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

type NewCallbackWebookMessage struct {
	webhookMessageBase
	Data NewCallbackWebhookData `json:"data"`
}

type NewCallbackWebhookData struct {
	User           string `json:"user" mapstructure:"user"`
	Host           string `json:"host" mapstructure:"host"`
	IPs            string `json:"ips" mapstructure:"ips"`
	Domain         string `json:"domain" mapstructure:"domain"`
	ExternalIP     string `json:"external_ip" mapstructure:"external_ip"`
	ProcessName    string `json:"process_name" mapstructure:"process_name"`
	PID            int    `json:"pid" mapstructure:"pid"`
	Os             string `json:"os" mapstructure:"os"`
	Architecture   string `json:"architecture" mapstructure:"architecture"`
	AgentType      string `json:"agent_type" mapstructure:"agent_type"`
	Description    string `json:"description" mapstructure:"description"`
	ExtraInfo      string `json:"extra_info" mapstructure:"extra_info"`
	SleepInfo      string `json:"sleep_info" mapstructure:"sleep_info"`
	DisplayID      int    `json:"display_id" mapstructure:"display_id"`
	ID             int    `json:"id" mapstructure:"id"`
	IntegrityLevel int    `json:"integrity_level" mapstructure:"integrity_level"`
}

// Register this method with rabbitmq so it can be called
func init() {
	AllWebhookData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         WEBHOOK_TYPE_NEW_CALLBACK,
		RabbitmqProcessingFunction: processNewCallbackWebhook,
	})
}
func processNewCallbackWebhook(input []byte) {
	inputStruct := NewCallbackWebookMessage{}
	if err := json.Unmarshal(input, &inputStruct); err != nil {
		logging.LogError(err, "Failed to process new callback webhook message")
	} else {
		// success, so do RPC calls to Mythic to get more context or send off webhook now
		for _, webhook := range AllWebhookData.GetAllNames() {
			if AllWebhookData.Get(webhook).GetWebhookDefinition().NewCallbackFunction != nil {
				AllWebhookData.Get(webhook).GetWebhookDefinition().NewCallbackFunction(inputStruct)
			}
		}
	}
}

package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackSearchCommandMessage struct {
	CallbackID                *int      `json:"callback_id,omitempty"`
	TaskID                    *int      `json:"task_id,omitempty"`
	SearchCommandNames        *[]string `json:"command_names,omitempty"`
	SearchSupportedUIFeatures *string   `json:"supported_ui_features,omitempty"`
	SearchScriptOnly          *bool     `json:"script_only,omitempty"`
	// this is an exact match search
	SearchAttributes map[string]interface{} `json:"params,omitempty"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCCallbackSearchCommandMessageResponse struct {
	Success  bool                                `json:"success"`
	Error    string                              `json:"error"`
	Commands []MythicRPCCommandSearchCommandData `json:"commands"`
}

type MythicRPCCallbackSearchCommandData struct {
	Name                string                 `json:"cmd"`
	Version             int                    `json:"version"`
	Attributes          map[string]interface{} `json:"attributes"`
	NeedsAdmin          bool                   `json:"needs_admin"`
	HelpCmd             string                 `json:"help_cmd"`
	Description         string                 `json:"description"`
	SupportedUiFeatures []string               `json:"supported_ui_features"`
	Author              string                 `json:"author"`
	ScriptOnly          bool                   `json:"script_only"`
}

func SendMythicRPCCallbackSearchCommand(input MythicRPCCallbackSearchCommandMessage) (*MythicRPCCallbackSearchCommandMessageResponse, error) {
	response := MythicRPCCallbackSearchCommandMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_SEARCH_COMMAND,
		input,
	); err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	} else if err := json.Unmarshal(responseBytes, &response); err != nil {
		logging.LogError(err, "Failed to parse SendMythicRPCCommandSearch response back to struct", "response", response)
		return nil, err
	} else {
		return &response, nil
	}
}

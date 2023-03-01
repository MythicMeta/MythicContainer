package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCommandSearchMessage struct {
	SearchCommandNames        *[]string `json:"command_names,omitempty"`
	SearchPayloadTypeName     string    `json:"payload_type_name"`
	SearchSupportedUIFeatures *string   `json:"supported_ui_features,omitempty"`
	SearchScriptOnly          *bool     `json:"script_only,omitempty"`
	SearchOs                  *string   `json:"os,omitempty"`
	// this is an exact match search
	SearchAttributes map[string]interface{} `json:"params,omitempty"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCCommandSearchMessageResponse struct {
	Success  bool                                `json:"success"`
	Error    string                              `json:"error"`
	Commands []MythicRPCCommandSearchCommandData `json:"commands"`
}

type MythicRPCCommandSearchCommandData struct {
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

func SendMythicRPCCommandSearch(input MythicRPCCommandSearchMessage) (*MythicRPCCommandSearchMessageResponse, error) {
	response := MythicRPCCommandSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_COMMAND_SEARCH,
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

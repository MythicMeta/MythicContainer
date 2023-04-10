package mythicrpc

import (
	"encoding/json"
	"time"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackSearchMessage struct {
	// AgentCallbackID (Required) - this is the UUID of the callback associated with this search.
	// This provides the necessary context to scope the search to the right operation.
	AgentCallbackID int `json:"agent_callback_id"`
	// SearchCallbackID (Optional) - if you know the real callback ID, you can search via that here.
	SearchCallbackID *int `json:"search_callback_id"`
	// SearchCallbackDisplayID (Optional) - if you know the display id for the callback (the one that shows up in the UI), then you can search via that here.
	SearchCallbackDisplayID *int `json:"search_callback_display_id"`
	// SearchCallbackUUID (Optional) - if you know the agent callback uuid for the callback, you can search for that here.
	SearchCallbackUUID *string `json:"search_callback_uuid"`
	// SearchCallbackUser (Optional) - if you know the user associated with the callback you want, supply that here.
	SearchCallbackUser *string `json:"user,omitempty"`
	// SearchCallbackHost (Optional) - if you know the hostname of the callback you want, supply that here.
	SearchCallbackHost *string `json:"host,omitempty"`
	// SearchCallbackPID (Optional) - if you know the PID of the callback you want, supply that here.
	SearchCallbackPID *int `json:"pid,omitempty"`
	// SearchCallbackExtraInfo (Optional) - if you know the extra info associated with a callback, supply that here.
	SearchCallbackExtraInfo *string `json:"extra_info,omitempty"`
	// SearchCallbackSleepInfo (Optional) - if you know the sleep information for a callback, supply that here.
	SearchCallbackSleepInfo *string `json:"sleep_info,omitempty"`
	// SearchCallbackIP (Optional) - if you know the IP address of the callback you want, supply that here
	SearchCallbackIP *string `json:"ip,omitempty"`
	// SearchCallbackExternalIP (Optional) - if you know the external IP address of the callback you want, supply that here.
	SearchCallbackExternalIP *string `json:"external_ip,omitempty"`
	// SearchCallbackIntegrityLevel (Optional) - if you know the integrity level of the callback you want, supply that here
	SearchCallbackIntegrityLevel *int `json:"integrity_level,omitempty"`
	// SearchCallbackOs (Optional) - if you know the detailed OS information for the callback you want, supply that here.
	// NOTE: This is NOT the "windows", "Linux", "macOS", etc piece you selected when building a payload.
	SearchCallbackOs *string `json:"os,omitempty"`
	// SearchCallbackDomain (Optional) - if you know the domain
	SearchCallbackDomain       *string `json:"domain,omitempty"`
	SearchCallbackArchitecture *string `json:"architecture,omitempty"`
	SearchCallbackDescription  *string `json:"description,omitempty"`
}
type MythicRPCCallbackSearchMessageResult struct {
	ID                  int       `mapstructure:"id" json:"id"`
	DisplayID           int       `mapstructure:"display_id" json:"display_id"`
	AgentCallbackID     string    `mapstructure:"agent_callback_id" json:"agent_callback_id"`
	InitCallback        time.Time `mapstructure:"init_callback" json:"init_callback"`
	LastCheckin         time.Time `mapstructure:"last_checkin" json:"last_checkin"`
	User                string    `mapstructure:"user" json:"user"`
	Host                string    `mapstructure:"host" json:"host"`
	PID                 int       `mapstructure:"pid" json:"pid"`
	Ip                  string    `mapstructure:"ip" json:"ip"`
	ExternalIp          string    `mapstructure:"external_ip" json:"external_ip"`
	ProcessName         string    `mapstructure:"process_name" json:"process_name"`
	Description         string    `mapstructure:"description" json:"description"`
	OperatorID          int       `mapstructure:"operator_id" json:"operator_id"`
	Active              bool      `mapstructure:"active" json:"active"`
	RegisteredPayloadID int       `mapstructure:"registered_payload_id" json:"registered_payload_id"`
	IntegrityLevel      int       `mapstructure:"integrity_level" json:"integrity_level"`
	Locked              bool      `mapstructure:"locked" json:"locked"`
	LockedOperatorID    int       `mapstructure:"locked_operator_id" json:"locked_operator_id"`
	OperationID         int       `mapstructure:"operation_id" json:"operation_id"`
	CryptoType          string    `mapstructure:"crypto_type" json:"crypto_type"`
	DecKey              *[]byte   `mapstructure:"dec_key" json:"dec_key"`
	EncKey              *[]byte   `mapstructure:"enc_key" json:"enc_key"`
	Os                  string    `mapstructure:"os" json:"os"`
	Architecture        string    `mapstructure:"architecture" json:"architecture"`
	Domain              string    `mapstructure:"domain" json:"domain"`
	ExtraInfo           string    `mapstructure:"extra_info" json:"extra_info"`
	SleepInfo           string    `mapstructure:"sleep_info" json:"sleep_info"`
	Timestamp           time.Time `mapstructure:"timestamp" json:"timestamp"`
}
type MythicRPCCallbackSearchMessageResponse struct {
	Success bool                                   `json:"success"`
	Error   string                                 `json:"error"`
	Results []MythicRPCCallbackSearchMessageResult `json:"results"`
}

func SendMythicRPCCallbackSearch(input MythicRPCCallbackSearchMessage) (*MythicRPCCallbackRemoveCommandMessageResponse, error) {
	response := MythicRPCCallbackRemoveCommandMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_SEARCH,
		input,
	); err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	} else if err := json.Unmarshal(responseBytes, &response); err != nil {
		logging.LogError(err, "Failed to parse response back to struct", "response", response)
		return nil, err
	} else {
		return &response, nil
	}
}

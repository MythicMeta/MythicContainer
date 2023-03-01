package mythicrpc

import (
	"encoding/json"
	"time"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackSearchMessage struct {
	AgentCallbackID              int     `json:"callback_id"` // required
	SearchCallbackID             *int    `json:"search_callback_id"`
	SearchCallbackDisplayID      *int    `json:"search_callback_display_id"`
	SearchCallbackUUID           *string `json:"search_callback_uuid"`
	SearchCallbackUser           *string `json:"user,omitempty"`
	SearchCallbackHost           *string `json:"host,omitempty"`
	SearchCallbackPID            *int    `json:"pid,omitempty"`
	SearchCallbackExtraInfo      *string `json:"extra_info,omitempty"`
	SearchCallbackSleepInfo      *string `json:"sleep_info,omitempty"`
	SearchCallbackIP             *string `json:"ip,omitempty"`
	SearchCallbackExternalIP     *string `json:"external_ip,omitempty"`
	SearchCallbackIntegrityLevel *int    `json:"integrity_level,omitempty"`
	SearchCallbackOs             *string `json:"os,omitempty"`
	SearchCallbackDomain         *string `json:"domain,omitempty"`
	SearchCallbackArchitecture   *string `json:"architecture,omitempty"`
	SearchCallbackDescription    *string `json:"description,omitempty"`
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

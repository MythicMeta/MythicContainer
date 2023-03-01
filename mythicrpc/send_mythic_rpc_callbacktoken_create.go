package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackTokenCreateMessage struct {
	TaskID         int                          `json:"task_id"` //required
	CallbackTokens []MythicRPCCallbackTokenData `json:"callbacktokens"`
}
type MythicRPCCallbackTokenCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type MythicRPCCallbackTokenData = agentMessagePostResponseCallbackTokens
type agentMessagePostResponseToken struct {
	Action             string `mapstructure:"action"`
	TokenID            int    `mapstructure:"token_id"`
	User               string `mapstructure:"user"`
	Groups             string `mapstructure:"groups"`
	Privileges         string `mapstructure:"privileges"`
	ThreadID           int    `mapstructure:"thread_id"`
	ProcessID          int    `mapstructure:"process_id"`
	SessionID          int    `mapstructure:"session_id"`
	LogonSID           string `mapstructure:"logon_sid"`
	IntegrityLevelSID  string `mapstructure:"integrity_level_sid"`
	Restricted         bool   `mapstructure:"restricted"`
	DefaultDacl        string `mapstructure:"default_dacl"`
	Handle             int    `mapstructure:"handle"`
	Capabilities       string `mapstructure:"capabilities"`
	AppContainerSID    string `mapstructure:"app_container_sid"`
	AppContainerNumber int    `mapstructure:"app_container_number"`
}
type agentMessagePostResponseCallbackTokens struct {
	Action  string  `mapstructure:"action"`
	Host    *string `mapstructure:"host,omitempty"`
	TokenId int     `mapstructure:"TokenId"`
	// optionally also provide all the token information
	TokenInfo *agentMessagePostResponseToken `mapstructure:"token"`
}

func SendMythicRPCCallbackTokenCreate(input MythicRPCCallbackTokenCreateMessage) (*MythicRPCCallbackTokenCreateMessageResponse, error) {
	response := MythicRPCCallbackTokenCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACKTOKEN_CREATE,
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

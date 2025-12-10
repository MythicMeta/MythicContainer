package mythicrpc

import (
	"encoding/json"
	"time"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackTokenSearchMessage struct {
	TaskID          *int    `json:"task_id"`
	CallbackID      *int    `json:"callback_id"`
	AgentCallbackID *string `json:"agent_callback_id"`
}
type MythicRPCCallbackTokenSearchMessageResponse struct {
	Success        bool                               `json:"success"`
	Error          string                             `json:"error"`
	CallbackTokens []MythicRPCCallbackSearchTokenData `json:"callbacktokens"`
}
type MythicRPCCallbackSearchTokenData struct {
	ID               int                               `mapstructure:"id" json:"id"`
	Token            MythicRPCCallbackTokenSearchToken `mapstructure:"token" json:"token"`
	CallbackID       int                               `mapstructure:"callback_id" json:"callback_id"`
	TaskID           int                               `mapstructure:"task_id" json:"task_id"`
	TimestampCreated time.Time                         `mapstructure:"timestamp_created" json:"timestamp_created"`
	Deleted          bool                              `mapstructure:"deleted" json:"deleted"`
	Host             string                            `mapstructure:"host" json:"host"`
}
type MythicRPCCallbackTokenSearchToken struct {
	// mythic supplied
	ID          int       `mapstructure:"id" json:"id"`
	TaskID      int       `mapstructure:"task_id" json:"task_id"`
	Deleted     bool      `mapstructure:"deleted" json:"deleted"`
	Host        string    `mapstructure:"host" json:"host"`
	Description string    `mapstructure:"description" json:"description"`
	OperationID int       `mapstructure:"operation_id" json:"operation_id"`
	Timestamp   time.Time `mapstructure:"timestamp" json:"timestamp"`
	// agent supplied
	TokenID            int64  `mapstructure:"token_id" json:"token_id"`
	User               string `mapstructure:"user" json:"user"`
	Groups             string `mapstructure:"groups" json:"groups"`
	Privileges         string `mapstructure:"privileges" json:"privileges"`
	ThreadID           int    `mapstructure:"thread_id" json:"thread_id"`
	ProcessID          int    `mapstructure:"process_id" json:"process_id"`
	SessionID          int    `mapstructure:"session_id" json:"session_id"`
	LogonSID           string `mapstructure:"logon_sid" json:"logon_sid"`
	IntegrityLevelSID  string `mapstructure:"integrity_level_sid" json:"integrity_level_sid"`
	AppContainerSID    string `mapstructure:"app_container_sid" json:"app_container_sid"`
	AppContainerNumber int    `mapstructure:"app_container_number" json:"app_container_number"`
	DefaultDacl        string `mapstructure:"default_dacl" json:"default_dacl"`
	Restricted         bool   `mapstructure:"restricted" json:"restricted"`
	Handle             int    `mapstructure:"handle" json:"handle"`
	Capabilities       string `mapstructure:"capabilities" json:"capabilities"`
}

func SendMythicRPCCallbackTokenSearch(input MythicRPCCallbackTokenSearchMessage) (*MythicRPCCallbackTokenSearchMessageResponse, error) {
	response := MythicRPCCallbackTokenSearchMessageResponse{}
	responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACKTOKEN_SEARCH,
		input,
	)
	if err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		logging.LogError(err, "Failed to parse response back to struct", "response", response)
		return nil, err
	}
	return &response, nil
}

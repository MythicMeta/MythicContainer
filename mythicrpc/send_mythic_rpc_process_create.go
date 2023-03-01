package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCProcessCreateMessage struct {
	TaskID    int                                 `json:"task_id"` //required
	Processes []MythicRPCProcessCreateProcessData `json:"processes"`
}
type MythicRPCProcessCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type MythicRPCProcessCreateProcessData = agentMessagePostResponseProcesses
type agentMessagePostResponseProcesses struct {
	Host                   *string                `mapstructure:"host,omitempty" json:"host,omitempty"`
	ProcessID              int                    `mapstructure:"process_id" json:"process_id"`
	ParentProcessID        int                    `mapstructure:"parent_process_id" json:"parent_process_id"`
	Architecture           string                 `mapstructure:"architecture" json:"architecture"`
	BinPath                string                 `mapstructure:"bin_path" json:"bin_path"`
	Name                   string                 `mapstructure:"name" json:"name"`
	User                   string                 `mapstructure:"user" json:"user"`
	CommandLine            string                 `mapstructure:"command_line" json:"command_line"`
	IntegrityLevel         int                    `mapstructure:"integrity_level" json:"integrity_level"`
	StartTime              int                    `mapstructure:"start_time" json:"start_time"`
	Description            string                 `mapstructure:"description" json:"description"`
	Signer                 string                 `mapstructure:"signer" json:"signer"`
	ProtectionProcessLevel int                    `mapstructure:"protected_process_level" json:"protected_process_level"`
	Other                  map[string]interface{} `json:"-" mapstructure:",remain"`
}

func SendMythicRPCProcessCreate(input MythicRPCProcessCreateMessage) (*MythicRPCProcessCreateMessageResponse, error) {
	response := MythicRPCProcessCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_PROCESS_CREATE,
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

package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCProcessSearchMessage struct {
	TaskID        int                               `json:"task_id"` //required
	SearchProcess MythicRPCProcessSearchProcessData `json:"process"`
}
type MythicRPCProcessSearchMessageResponse struct {
	Success   bool                                `json:"success"`
	Error     string                              `json:"error"`
	Processes []MythicRPCProcessSearchProcessData `json:"processes"`
}
type MythicRPCProcessSearchProcessData struct {
	Host            *string `json:"host" `              // optional
	ProcessID       *int    `json:"process_id" `        // optional
	Architecture    *string `json:"architecture"`       // optional
	ParentProcessID *int    `json:"parent_process_id" ` // optional
	BinPath         *string `json:"bin_path" `          // optional
	Name            *string `json:"name" `              // optional
	User            *string `json:"user" `              // optional
	CommandLine     *string `json:"command_line" `      // optional
	IntegrityLevel  *int    `json:"integrity_level" `   // optional
	Description     *string `json:"description" `       // optional
	Signer          *string `json:"signer"`             // optional
}

func SendMythicRPCProcessSearch(input MythicRPCProcessSearchMessage) (*MythicRPCProcessSearchMessageResponse, error) {
	response := MythicRPCProcessSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_PROCESS_SEARCH,
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

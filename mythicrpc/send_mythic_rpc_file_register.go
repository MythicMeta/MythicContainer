package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCFileRegisterMessage struct {
	Filename         string `json:"filename"`
	Comment          string `json:"comment"`
	OperationID      int    `json:"operation_id"`
	OperatorID       int    `json:"operator_id"`
	DeleteAfterFetch bool   `json:"delete_after_fetch"`
}
type MythicRPCFileRegisterMessageResponse struct {
	Success     bool   `json:"success"`
	Error       string `json:"error"`
	AgentFileId string `json:"agent_file_id"`
}

func SendMythicRPCFileRegister(input MythicRPCFileRegisterMessage) (*MythicRPCFileRegisterMessageResponse, error) {
	response := MythicRPCFileRegisterMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_FILE_REGISTER,
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

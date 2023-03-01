package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCredentialCreateMessage struct {
	TaskID      int                                       `json:"task_id"` //required
	Credentials []MythicRPCCredentialCreateCredentialData `json:"credentials"`
}
type MythicRPCCredentialCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type MythicRPCCredentialCreateCredentialData = agentMessagePostResponseCredentials
type agentMessagePostResponseCredentials struct {
	CredentialType string `json:"credential_type" mapstructure:"credential_type"`
	Realm          string `json:"realm" mapstructure:"realm"`
	Account        string `json:"account" mapstructure:"account"`
	Credential     string `json:"credential" mapstructure:"credential"`
	Comment        string `json:"comment" mapstructure:"comment"`
	ExtraData      string `json:"metadata" mapstructure:"metadata"`
}

func SendMythicRPCCredentialCreate(input MythicRPCCommandSearchMessage) (*MythicRPCCredentialCreateMessageResponse, error) {
	response := MythicRPCCredentialCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CREDENTIAL_CREATE,
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

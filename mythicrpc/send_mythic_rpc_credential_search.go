package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCredentialSearchMessage struct {
	TaskID            int                                     `json:"task_id"` //required
	SearchCredentials MythicRPCCredentialSearchCredentialData `json:"credentials"`
}
type MythicRPCCredentialSearchMessageResponse struct {
	Success     bool                                      `json:"success"`
	Error       string                                    `json:"error"`
	Credentials []MythicRPCCredentialSearchCredentialData `json:"credentials"`
}
type MythicRPCCredentialSearchCredentialData struct {
	Type       *string `json:"type" `      // optional
	Account    *string `json:"account" `   // optional
	Realm      *string `json:"realm" `     // optional
	Credential *string `json:"credential"` // optional
	Comment    *string `json:"comment"`    // optional
	Metadata   *string `json:"metadata"`   // optional
	Task_ID    int     `json:"task_id"`    // optional
}

func SendMythicRPCCredentialSearch(input MythicRPCCredentialSearchMessage) (*MythicRPCCredentialSearchMessageResponse, error) {
	response := MythicRPCCredentialSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CREDENTIAL_SEARCH,
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

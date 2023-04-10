package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCAgentstorageCreateMessage struct {
	// UniqueID (Required) - a unique identifier for this entry provided by you, the dev
	UniqueID string `json:"unique_id"`
	// DataToStore (Required) - the data you want to store as bytes
	DataToStore []byte `json:"data"`
}
type MythicRPCAgentstorageCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// SendMythicRPCAgentStorageCreate - Create a new entry in the agentstorage table within Mythic.
// This can be used to store arbitrary data that the agent/c2 profile might need later on and used a way to share data.
func SendMythicRPCAgentStorageCreate(input MythicRPCAgentstorageCreateMessage) (*MythicRPCAgentstorageCreateMessageResponse, error) {
	response := MythicRPCAgentstorageCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_AGENTSTORAGE_CREATE,
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

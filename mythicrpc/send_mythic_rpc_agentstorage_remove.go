package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCAgentstorageRemoveMessage struct {
	// UniqueID (Required) - The unique identifier to search for in the database to remove
	UniqueID string `json:"unique_id"`
}
type MythicRPCAgentstorageRemoveMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// SendMythicRPCAgentStorageRemove - Remove a specific entry from the agentstorage table within Mythic.
func SendMythicRPCAgentStorageRemove(input MythicRPCAgentstorageRemoveMessage) (*MythicRPCAgentstorageRemoveMessageResponse, error) {
	response := MythicRPCAgentstorageRemoveMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_AGENTSTORAGE_REMOVE,
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

package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCAgentstorageSearchMessage struct {
	// SearchUniqueID (Required) - The unique identifier you supplied when creating the data that you're searching for
	SearchUniqueID string `json:"unique_id"` // required
}
type MythicRPCAgentstorageSearchMessageResponse struct {
	Success              bool                                `json:"success"`
	Error                string                              `json:"error"`
	AgentStorageMessages []MythicRPCAgentstorageSearchResult `json:"agentstorage_messages"`
}
type MythicRPCAgentstorageSearchResult struct {
	UniqueID string `json:"unique_id"`
	Data     []byte `json:"data"`
}

// SendMythicRPCAgentStorageSearch - Search for a specific entry within the agentstorage table and fetch the results.
func SendMythicRPCAgentStorageSearch(input MythicRPCAgentstorageSearchMessage) (*MythicRPCAgentstorageSearchMessageResponse, error) {
	response := MythicRPCAgentstorageSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_AGENTSTORAGE_SEARCH,
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

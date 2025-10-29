package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackEdgeRemoveMessage struct {
	SourceCallbackID      int    `json:"source_callback_id"`
	DestinationCallbackID int    `json:"destination_callback_id"`
	C2ProfileName         string `json:"c2_profile_name"`
}
type MythicRPCCallbackEdgeRemoveMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendMythicRPCCallbackEdgeRemove(input MythicRPCCallbackEdgeRemoveMessage) (*MythicRPCCallbackEdgeRemoveMessageResponse, error) {
	response := MythicRPCCallbackEdgeRemoveMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_EDGE_REMOVE,
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

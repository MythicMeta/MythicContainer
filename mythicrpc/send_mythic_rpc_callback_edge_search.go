package mythicrpc

import (
	"encoding/json"
	"time"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackEdgeSearchMessage struct {
	AgentCallbackUUID     string  `json:"agent_callback_id"`
	AgentCallbackID       int     `json:"callback_id"`
	SearchC2ProfileName   *string `json:"search_c2_profile_name"`
	SearchActiveEdgesOnly *bool   `json:"search_active_edges_only"`
}
type MythicRPCCallbackEdgeSearchMessageResult struct {
	ID             int                                  `mapstructure:"id" json:"id"`
	StartTimestamp time.Time                            `mapstructure:"start_timestamp" json:"start_timestamp"`
	EndTimestamp   time.Time                            `mapstructure:"end_timestamp" json:"end_timestamp"`
	Source         MythicRPCCallbackSearchMessageResult `mapstructure:"source" json:"source"`
	Destination    MythicRPCCallbackSearchMessageResult `mapstructure:"destination" json:"destination"`
	C2Profile      string                               `mapstructure:"c2profile" json:"c2profile"`
}
type MythicRPCCallbackEdgeSearchMessageResponse struct {
	Success bool                                       `json:"success"`
	Error   string                                     `json:"error"`
	Results []MythicRPCCallbackEdgeSearchMessageResult `json:"results"`
}

func SendMythicRPCCallbackEdgeSearch(input MythicRPCCallbackEdgeSearchMessage) (*MythicRPCCallbackEdgeSearchMessageResponse, error) {
	response := MythicRPCCallbackEdgeSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_EDGE_SEARCH,
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

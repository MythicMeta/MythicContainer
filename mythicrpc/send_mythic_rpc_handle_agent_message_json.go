package mythicrpc

import (
	"context"
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCHandleAgentJsonMessage struct {
	// CallbackID or AgentCallbackID is required
	CallbackID      *int    `json:"callback_id"`
	AgentCallbackID *string `json:"agent_callback_id"`
	// AgentMessage is the full agent JSON message including an action keyword
	AgentMessage      map[string]interface{} `json:"agent_message"`
	UpdateCheckinTime bool                   `json:"update_checkin_time"`
}

type MythicRPCHandleAgentMessageJsonMessage struct {
	// CallbackID or AgentCallbackID is required
	CallbackID      int    `json:"callback_id,omitempty"`
	AgentCallbackID string `json:"agent_callback_id,omitempty"`
	// AgentMessage is the full agent JSON message including an action keyword
	AgentMessage      map[string]interface{} `json:"agent_message"`
	UpdateCheckinTime bool                   `json:"update_checkin_time"`
}

type MythicRPCHandleAgentJsonMessageResponse struct {
	Success       bool                   `json:"success"`
	Error         string                 `json:"error"`
	AgentResponse map[string]interface{} `json:"agent_response"`
}

type MythicRPCHandleAgentMessageJsonMessageResponse = MythicRPCHandleAgentJsonMessageResponse

func SendMythicRPCHandleAgentJson(ctx context.Context, input MythicRPCHandleAgentJsonMessage) (*MythicRPCHandleAgentJsonMessageResponse, error) {
	response := MythicRPCHandleAgentJsonMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessageWithContext(
		ctx,
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_HANDLE_AGENT_JSON,
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

func SendMythicRPCHandleAgentMessageJson(ctx context.Context, input MythicRPCHandleAgentMessageJsonMessage) (*MythicRPCHandleAgentMessageJsonMessageResponse, error) {
	response := MythicRPCHandleAgentMessageJsonMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessageWithContext(
		ctx,
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_HANDLE_AGENT_JSON,
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

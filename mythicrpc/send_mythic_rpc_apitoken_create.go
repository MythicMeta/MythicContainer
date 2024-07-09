package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

// MythicRPCAPITokenCreateMessage needs at least one parameter to generate an appropriate apitoken, the rest are unnecessary
type MythicRPCAPITokenCreateMessage struct {
	AgentTaskID     *string `json:"agent_task_id"`
	AgentCallbackID *string `json:"agent_callback_id"`
	PayloadUUID     *string `json:"payload_uuid"`
	OperationID     *int    `json:"operation_id"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCAPITokenCreateMessageResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	APIToken string `json:"apitoken"`
}

func SendMythicRPCAPITokenCreate(input MythicRPCAPITokenCreateMessage) (*MythicRPCAPITokenCreateMessageResponse, error) {
	response := MythicRPCAPITokenCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_APITOKEN_CREATE,
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

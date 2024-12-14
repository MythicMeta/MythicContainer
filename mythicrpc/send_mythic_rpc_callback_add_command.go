package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackAddCommandMessage struct {
	// TaskID - What task is trying to add commands. This will add commands to the callback associated with this task.
	TaskID int `json:"task_id"` // required
	// AgentCallbackID - Agent Callback UUID of callback to add commands to if TaskID isn't specified
	AgentCallbackID string `json:"agent_callback_id"`
	// PayloadType - The payload type of the associated commands to load if they're for a payload type (or command augment) container other than the one for the callback
	PayloadType string `json:"payload_type"`
	// Commands (Required) - The names of the commands you want to add. If they're already added, then they are skipped.
	Commands []string `json:"commands"` // required
	// CallbackIDs allows you to specify a list of CallbackID values to add commands to multiple callbacks at once
	CallbackIDs []int `json:"callback_ids"`
}
type MythicRPCCallbackAddCommandMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// SendMythicRPCCallbackAddCommand - Register new commands as being "loaded" into the current callback. This makes them
// available for tasking through the UI.
func SendMythicRPCCallbackAddCommand(input MythicRPCCallbackAddCommandMessage) (*MythicRPCCallbackAddCommandMessageResponse, error) {
	response := MythicRPCCallbackAddCommandMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_ADD_COMMAND,
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

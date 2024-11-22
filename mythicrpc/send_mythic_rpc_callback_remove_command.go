package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackRemoveCommandMessage struct {
	// TaskID (Required) - The task id that's going to remove commands from the associated callback.
	TaskID int `json:"task_id"` // required
	// AgentCallbackID - Agent Callback UUID of callback to add commands to if TaskID isn't specified
	AgentCallbackID string `json:"agent_callback_id"`
	// PayloadType - The payload type of the associated commands to load if they're for a payload type (or command augment) container other than the one for the callback
	PayloadType string `json:"payload_type"`
	// Commands (Required) - The list of command names to be removed from the callback. If the command isn't loaded
	// within the callback, then it's skipped
	Commands []string `json:"commands"` // required
}
type MythicRPCCallbackRemoveCommandMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// SendMythicRPCCallbackRemoveCommand - Remove commands from a certain callback. This is helpful if you want to
// unload certain functionality that might have been temporarily loaded in the first place.
func SendMythicRPCCallbackRemoveCommand(input MythicRPCCallbackRemoveCommandMessage) (*MythicRPCCallbackRemoveCommandMessageResponse, error) {
	response := MythicRPCCallbackRemoveCommandMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_REMOVE_COMMAND,
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

package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCallbackDisplayToRealIdSearchMessage struct {
	// CallbackDisplayID (Required) - The display id that the operator sees for a callback (i.e. 1, 2, 3, etc).
	// The display ID is always incremental within each operation. So, each operation will have a callback1 for example.
	CallbackDisplayID int `json:"callback_display_id"`
	// OperationName (Optional) - The name of the operation associated with the callback.
	// Either OperationName or OperationID must be supplied to give context for the CallbackDisplayID
	OperationName *string `json:"operation_name"`
	// OperationID (Opational) - The ID of the operation associated with the callback.
	// Either OperationName or OperationID must be supplied to give context for the CallbackDisplayID
	OperationID *int `json:"operation_id"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCCallbackDisplayToRealIdSearchMessageResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	CallbackID int    `json:"callback_id"`
}

// SendMythicRPCCallbackDisplayToRealIdSearch - Convert a generic CallbackDisplayID to a unique CallbackID for use with other RPC calls.
func SendMythicRPCCallbackDisplayToRealIdSearch(input MythicRPCCallbackDisplayToRealIdSearchMessage) (*MythicRPCCallbackDisplayToRealIdSearchMessageResponse, error) {
	response := MythicRPCCallbackDisplayToRealIdSearchMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_CALLBACK_DISPLAY_TO_REAL_ID_SEARCH,
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

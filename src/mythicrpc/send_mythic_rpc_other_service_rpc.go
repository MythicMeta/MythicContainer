package mythicrpc

import (
	"encoding/json"
	"fmt"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCOtherServiceRPCMessage struct {
	ServiceName                 string                 `json:"service_name"` //required
	ServiceRPCFunction          string                 `json:"service_function"`
	ServiceRPCFunctionArguments map[string]interface{} `json:"service_arguments"`
}
type MythicRPCOtherServiceRPCMessageResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Result  map[string]interface{} `json:"result"`
}

func getMythicRPCOtherServiceRPCRoutingKey(service string) string {
	return fmt.Sprintf("%s_%s", service, rabbitmq.MYTHIC_RPC_OTHER_SERVICES_RPC)
}
func SendMythicRPCOtherServiceRPC(input MythicRPCOtherServiceRPCMessage) (*MythicRPCOtherServiceRPCMessageResponse, error) {
	response := MythicRPCOtherServiceRPCMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		getMythicRPCOtherServiceRPCRoutingKey(input.ServiceName),
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

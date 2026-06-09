package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	c2structs "github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_DYNAMIC_QUERY_C2_PARAMETER_FUNCTION,
		RabbitmqProcessingFunction: processC2RPCC2ParameterDynamicQueryFunctionMessages,
	})
}

func processC2RPCC2ParameterDynamicQueryFunctionMessages(ctx context.Context, msg []byte) interface{} {
	incomingMessage := c2structs.C2RPCDynamicQueryC2ParameterFunctionMessage{}
	response := c2structs.C2RPCDynamicQueryC2ParameterFunctionMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Error = "Failed to unmarshal JSON message into structs"
		return response
	}

	for _, param := range c2structs.AllC2Data.Get(incomingMessage.C2Profile).GetParameters() {
		if param.Name == incomingMessage.ParameterName {
			if param.DynamicQueryFunction != nil {
				response = param.DynamicQueryFunction(ctx, incomingMessage)
				response.Success = true
				return response
			}

			response.Choices = []string{}
			response.Error = "Function was nil"
			return response
		}
	}
	response.Error = fmt.Sprintf("Failed to find parameter %s", incomingMessage.ParameterName)
	return response
}

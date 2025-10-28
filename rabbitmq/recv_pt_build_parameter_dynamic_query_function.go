package rabbitmq

import (
	"encoding/json"
	"fmt"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         PT_RPC_DYNAMIC_QUERY_BUILD_PARAMETER_FUNCTION,
		RabbitmqProcessingFunction: processPtRPCBuildParameterDynamicQueryFunctionMessages,
	})
}

func processPtRPCBuildParameterDynamicQueryFunctionMessages(msg []byte) interface{} {
	incomingMessage := agentstructs.PTRPCDynamicQueryBuildParameterFunctionMessage{}
	response := agentstructs.PTRPCDynamicQueryBuildParameterFunctionMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Error = "Failed to unmarshal JSON message into structs"
		return response
	} else {
		for _, param := range agentstructs.AllPayloadData.Get(incomingMessage.PayloadType).GetBuildParameters() {
			if param.Name == incomingMessage.ParameterName {
				if param.DynamicQueryFunction != nil {
					response = param.DynamicQueryFunction(incomingMessage)
					response.Success = true
					return response
				} else {
					response.Choices = []string{}
					response.Error = "Function was nil"
					return response
				}
			}
		}
		response.Error = fmt.Sprintf("Failed to find parameter %s", incomingMessage.ParameterName)
		return response
	}
}

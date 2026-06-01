package rabbitmq

import (
	"context"
	"encoding/json"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         MYTHIC_RPC_OTHER_SERVICES_RPC,
		RabbitmqProcessingFunction: processC2OtherServiceRPC,
	})
	agentstructs.AllPayloadData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         MYTHIC_RPC_OTHER_SERVICES_RPC,
		RabbitmqProcessingFunction: processPTOtherServiceRPC,
	})
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processC2OtherServiceRPC(ctx context.Context, msg []byte) interface{} {
	input := c2structs.C2RPCOtherServiceRPCMessage{}
	responseMsg := c2structs.C2RPCOtherServiceRPCMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return C2OtherServiceRPC(ctx, input)
	}
	return responseMsg
}

func C2OtherServiceRPC(ctx context.Context, input c2structs.C2RPCOtherServiceRPCMessage) c2structs.C2RPCOtherServiceRPCMessageResponse {
	responseMsg := c2structs.C2RPCOtherServiceRPCMessageResponse{
		Success: false,
		Error:   "Failed to find function",
	}
	if c2structs.AllC2Data.Get(input.ServiceName).GetC2Definition().CustomRPCFunctions != nil {
		for name, _ := range c2structs.AllC2Data.Get(input.ServiceName).GetC2Definition().CustomRPCFunctions {
			if name == input.ServiceRPCFunction {
				responseMsg = c2structs.AllC2Data.Get(input.ServiceName).GetC2Definition().CustomRPCFunctions[input.ServiceRPCFunction](ctx, input)
			}
		}
	}
	if responseMsg.RestartInternalServer {
		go restartC2Server(input.ServiceName)
	}
	return responseMsg
}

func processPTOtherServiceRPC(ctx context.Context, msg []byte) interface{} {
	input := agentstructs.PTRPCOtherServiceRPCMessage{}
	responseMsg := agentstructs.PTRPCOtherServiceRPCMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return PTOtherServiceRPC(ctx, input)
	}
	return responseMsg
}

func PTOtherServiceRPC(ctx context.Context, input agentstructs.PTRPCOtherServiceRPCMessage) agentstructs.PTRPCOtherServiceRPCMessageResponse {
	response := agentstructs.PTRPCOtherServiceRPCMessageResponse{
		Success: false,
		Error:   "Failed to find function",
	}
	if agentstructs.AllPayloadData.Get(input.Name).GetPayloadDefinition().CustomRPCFunctions != nil {
		for name, _ := range agentstructs.AllPayloadData.Get(input.Name).GetPayloadDefinition().CustomRPCFunctions {
			if name == input.RPCFunction {
				return agentstructs.AllPayloadData.Get(input.Name).GetPayloadDefinition().CustomRPCFunctions[input.RPCFunction](ctx, input)
			}
		}
	}
	return response
}

package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_GET_IOC_ROUTING_KEY,
		RabbitmqProcessingFunction: processC2RPCGetIOC,
	})
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processC2RPCGetIOC(msg []byte) interface{} {
	input := c2structs.C2GetIOCMessage{}
	responseMsg := c2structs.C2GetIOCMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do C2RPCGetIOC
		return C2RPCGetIOC(input)
	}
	return responseMsg
}

func C2RPCGetIOC(input c2structs.C2GetIOCMessage) c2structs.C2GetIOCMessageResponse {
	response := c2structs.C2GetIOCMessageResponse{
		Success: true,
		Error:   "No IOCs configured",
	}
	if c2structs.AllC2Data.Get(input.Name).GetC2Definition().GetIOCFunction != nil {
		response = c2structs.AllC2Data.Get(input.Name).GetC2Definition().GetIOCFunction(input)
	}
	return response
}

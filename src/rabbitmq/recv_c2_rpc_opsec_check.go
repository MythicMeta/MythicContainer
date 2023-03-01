package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/c2_structs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(c2structs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_OPSEC_CHECKS_ROUTING_KEY,
		RabbitmqProcessingFunction: processC2RPCOpsecCheck,
	})
}

func processC2RPCOpsecCheck(msg []byte) interface{} {
	input := c2structs.C2OPSECMessage{}
	responseMsg := c2structs.C2OPSECMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return C2RPCOpsecCheck(input)
	}

	return responseMsg
}

func C2RPCOpsecCheck(input c2structs.C2OPSECMessage) c2structs.C2OPSECMessageResponse {
	responseMsg := c2structs.C2OPSECMessageResponse{
		Success: true,
		Error:   "No OPSEC Check performed - passing by default",
	}
	if c2structs.AllC2Data.Get(input.Name).GetC2Definition().OPSECCheckFunction != nil {
		responseMsg = c2structs.AllC2Data.Get(input.Name).GetC2Definition().OPSECCheckFunction(input)
	}
	return responseMsg
}

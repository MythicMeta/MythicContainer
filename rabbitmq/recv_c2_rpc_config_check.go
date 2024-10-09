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
		RabbitmqRoutingKey:         C2_RPC_CONFIG_CHECK_ROUTING_KEY,
		RabbitmqProcessingFunction: processC2RPCConfigCheck,
	})
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processC2RPCConfigCheck(msg []byte) interface{} {
	input := c2structs.C2ConfigCheckMessage{}
	responseMsg := c2structs.C2ConfigCheckMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return C2RPCConfigCheck(input)
	}
	return responseMsg
}

func C2RPCConfigCheck(input c2structs.C2ConfigCheckMessage) c2structs.C2ConfigCheckMessageResponse {
	responseMsg := c2structs.C2ConfigCheckMessageResponse{
		Success: true,
		Error:   "No Config Check performed - passing by default",
	}
	c2Mutex.Lock()
	if c2structs.AllC2Data.Get(input.Name).GetC2Definition().ConfigCheckFunction != nil {
		responseMsg = c2structs.AllC2Data.Get(input.Name).GetC2Definition().ConfigCheckFunction(input)
	}
	c2Mutex.Unlock()
	if responseMsg.RestartInternalServer {
		go restartC2Server(input.Name)
	}
	return responseMsg
}

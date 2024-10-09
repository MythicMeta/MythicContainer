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
		RabbitmqRoutingKey:         C2_RPC_SAMPLE_MESSAGE_ROUTING_KEY,
		RabbitmqProcessingFunction: processC2RPCSampleMessage,
	})
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processC2RPCSampleMessage(msg []byte) interface{} {
	input := c2structs.C2SampleMessageMessage{}
	responseMsg := c2structs.C2SampleMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do C2RPCGetIOC
		return C2RPCSampleMessage(input)
	}
	return responseMsg
}

func C2RPCSampleMessage(input c2structs.C2SampleMessageMessage) c2structs.C2SampleMessageResponse {
	responseMsg := c2structs.C2SampleMessageResponse{
		Success: true,
		Error:   "No Sample Message configured",
	}
	c2Mutex.Lock()
	if c2structs.AllC2Data.Get(input.Name).GetC2Definition().SampleMessageFunction != nil {
		responseMsg = c2structs.AllC2Data.Get(input.Name).GetC2Definition().SampleMessageFunction(input)
	}
	c2Mutex.Unlock()
	if responseMsg.RestartInternalServer {
		go restartC2Server(input.Name)
	}
	return responseMsg
}

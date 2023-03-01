package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/c2_structs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(c2structs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_RESYNC_ROUTING_KEY,
		RabbitmqProcessingFunction: processC2RPCReSync,
	})
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processC2RPCReSync(msg []byte) interface{} {
	input := c2structs.C2RPCReSyncMessage{}
	responseMsg := c2structs.C2RPCReSyncMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return C2RPCReSync(input)
	}
	return responseMsg
}

func C2RPCReSync(input c2structs.C2RPCReSyncMessage) c2structs.C2RPCReSyncMessageResponse {
	response := c2structs.C2RPCReSyncMessageResponse{
		Success: true,
		Error:   "",
	}
	SyncAllC2Data(&input.Name)
	return response
}

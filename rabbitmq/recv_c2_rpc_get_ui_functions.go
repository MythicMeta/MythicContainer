package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/c2_structs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(c2structs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_GET_AVAILABLE_UI_FUNCTIONS,
		RabbitmqProcessingFunction: processC2RPCGetUiFunctions,
	})
}

func processC2RPCGetUiFunctions(msg []byte) interface{} {
	input := c2structs.C2GetUiFunctionsMessage{}
	responseMsg := c2structs.C2GetUiFunctionsMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return C2RPCGetUiFunctions(input)
	}
	return responseMsg
}

func C2RPCGetUiFunctions(input c2structs.C2GetUiFunctionsMessage) c2structs.C2GetUiFunctionsMessageResponse {
	responseMsg := c2structs.C2GetUiFunctionsMessageResponse{
		Success: false,
		Error:   "Not implemented, not getting debug output",
	}
	return responseMsg
}

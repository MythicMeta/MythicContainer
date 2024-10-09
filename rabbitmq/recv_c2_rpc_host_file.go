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
		RabbitmqRoutingKey:         C2_RPC_HOST_FILE,
		RabbitmqProcessingFunction: processC2RPCHostFile,
	})
}

func processC2RPCHostFile(msg []byte) interface{} {
	input := c2structs.C2HostFileMessage{}
	responseMsg := c2structs.C2HostFileMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return C2RPCHostFile(input)
	}
	return responseMsg
}

func C2RPCHostFile(input c2structs.C2HostFileMessage) c2structs.C2HostFileMessageResponse {
	responseMsg := c2structs.C2HostFileMessageResponse{
		Success: false,
		Error:   "Not implemented, not hosting a file",
	}
	c2Mutex.Lock()
	if c2structs.AllC2Data.Get(input.Name).GetC2Definition().HostFileFunction != nil {
		responseMsg = c2structs.AllC2Data.Get(input.Name).GetC2Definition().HostFileFunction(input)
	}
	c2Mutex.Unlock()
	if responseMsg.RestartInternalServer {
		go restartC2Server(input.Name)
	}
	return responseMsg
}

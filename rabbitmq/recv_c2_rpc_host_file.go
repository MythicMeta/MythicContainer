package rabbitmq

import (
	"context"
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

func processC2RPCHostFile(ctx context.Context, msg []byte) interface{} {
	input := c2structs.C2HostFilesMessage{}
	responseMsg := c2structs.C2HostFilesMessageResponse{}
	err := json.Unmarshal(msg, &input)
	if err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
		return responseMsg
	}
	return C2RPCHostFile(ctx, input)
}

func C2RPCHostFile(ctx context.Context, input c2structs.C2HostFilesMessage) c2structs.C2HostFilesMessageResponse {
	responseMsg := c2structs.C2HostFilesMessageResponse{
		Success: false,
		Error:   "Not implemented, not hosting a file",
		Results: []c2structs.C2HostFileMessageResponse{},
	}
	for _, file := range input.Files {
		responseMsg.Results = append(responseMsg.Results, c2structs.C2HostFileMessageResponse{
			Success:     false,
			Error:       "Not implemented, not hosting a file",
			AgentFileID: file.AgentFileID,
			HostURL:     file.HostURL,
		})
	}
	c2Mutex.Lock()
	if c2structs.AllC2Data.Get(input.Name).GetC2Definition().HostFileFunction != nil {
		responseMsg = c2structs.AllC2Data.Get(input.Name).GetC2Definition().HostFileFunction(ctx, input)
	}
	c2Mutex.Unlock()
	if responseMsg.RestartInternalServer {
		go restartC2Server(ctx, input.Name)
	}
	return responseMsg
}

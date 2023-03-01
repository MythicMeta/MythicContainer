package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"os"
	"path/filepath"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(c2structs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_LIST_FILE,
		RabbitmqProcessingFunction: processC2RPCListFile,
	})
}

func processC2RPCListFile(msg []byte) interface{} {
	input := c2structs.C2RPCListFileMessage{}
	responseMsg := c2structs.C2RPCListFileMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return C2RPCListFile(input)
	}
	return responseMsg
}

func C2RPCListFile(input c2structs.C2RPCListFileMessage) c2structs.C2RPCListFileMessageResponse {
	responseMsg := c2structs.C2RPCListFileMessageResponse{
		Success: false,
	}
	if path, err := filepath.Abs(c2structs.AllC2Data.Get(input.Name).GetC2ServerFolderPath()); err != nil {
		logging.LogError(err, "Failed to get c2 server folder path")
		responseMsg.Error = err.Error()
		return responseMsg
	} else if entries, err := os.ReadDir(path); err != nil {
		logging.LogError(err, "Failed to list out contents of server folder path")
		responseMsg.Error = err.Error()
		return responseMsg
	} else {
		logging.LogInfo("getting file list", "folder path", c2structs.AllC2Data.Get(input.Name).GetC2ServerFolderPath())
		for _, entry := range entries {
			if !entry.IsDir() {
				responseMsg.Files = append(responseMsg.Files, entry.Name())
			}
		}
		responseMsg.Success = true
		return responseMsg
	}
}

package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"os"
	"path/filepath"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(c2structs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_REMOVE_FILE,
		RabbitmqProcessingFunction: processC2RPCRemoveFile,
	})
}

func processC2RPCRemoveFile(msg []byte) interface{} {
	input := c2structs.C2RPCRemoveFileMessage{}
	responseMsg := c2structs.C2RPCRemoveFileMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return C2RPCRemoveFile(input)
	}
	return responseMsg
}

func C2RPCRemoveFile(input c2structs.C2RPCRemoveFileMessage) c2structs.C2RPCRemoveFileMessageResponse {
	responseMsg := c2structs.C2RPCRemoveFileMessageResponse{
		Success: false,
	}
	if filePath, err := filepath.Abs(filepath.Join(c2structs.AllC2Data.Get(input.Name).GetC2ServerFolderPath(), input.Filename)); err != nil {
		logging.LogError(err, "Failed to get absolute filepath for file to remove")
		responseMsg.Error = fmt.Sprintf("Failed to locate file: %s\n", err.Error())
	} else if err := os.Remove(filePath); err != nil {
		responseMsg.Error = fmt.Sprintf("Failed to open file: %s", err.Error())
	} else {
		responseMsg.Success = true
	}
	return responseMsg
}

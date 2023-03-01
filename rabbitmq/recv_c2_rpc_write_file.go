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
		RabbitmqRoutingKey:         C2_RPC_WRITE_FILE,
		RabbitmqProcessingFunction: processC2RPCWriteFile,
	})
}

func processC2RPCWriteFile(msg []byte) interface{} {
	input := c2structs.C2RPCWriteFileMessage{}
	responseMsg := c2structs.C2RPCWriteFileMessageResponse{
		Success: false,
		Error:   "Not implemented, not getting debug output",
	}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return C2RPCWriteFile(input)
	}
	return responseMsg
}

func C2RPCWriteFile(input c2structs.C2RPCWriteFileMessage) c2structs.C2RPCWriteFileMessageResponse {
	responseMsg := c2structs.C2RPCWriteFileMessageResponse{
		Success: false,
	}
	if filePath, err := filepath.Abs(filepath.Join(c2structs.AllC2Data.Get(input.Name).GetC2ServerFolderPath(), input.Filename)); err != nil {
		logging.LogError(err, "Failed to get absolute filepath for file to get")
		responseMsg.Error = fmt.Sprintf("Failed to locate file: %s\n", err.Error())
	} else if err := os.WriteFile(filePath, input.Contents, 0644); err != nil {
		responseMsg.Error = fmt.Sprintf("Failed to open file: %s", err.Error())
	} else {
		responseMsg.Success = true
		responseMsg.Message = "Successfully wrote file"
	}
	return responseMsg
}

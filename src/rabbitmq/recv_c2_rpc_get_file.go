package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(c2structs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_GET_FILE,
		RabbitmqProcessingFunction: processC2RPCGetFile,
	})
}

func processC2RPCGetFile(msg []byte) interface{} {
	input := c2structs.C2RPCGetFileMessage{}
	responseMsg := c2structs.C2RPCGetFileMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return C2RPCGetFile(input)
	}
	return responseMsg
}

func C2RPCGetFile(input c2structs.C2RPCGetFileMessage) c2structs.C2RPCGetFileMessageResponse {
	responseMsg := c2structs.C2RPCGetFileMessageResponse{
		Success: false,
	}
	if filePath, err := filepath.Abs(filepath.Join(c2structs.AllC2Data.Get(input.Name).GetC2ServerFolderPath(), input.Filename)); err != nil {
		logging.LogError(err, "Failed to get absolute filepath for file to get")
		responseMsg.Error = fmt.Sprintf("Failed to locate file: %s\n", err.Error())
	} else if file, err := os.Open(filePath); err != nil {
		responseMsg.Error = fmt.Sprintf("Failed to open file: %s", err.Error())
	} else if contents, err := ioutil.ReadAll(file); err != nil {
		responseMsg.Error = fmt.Sprintf("Failed to read file: %s", err.Error())
	} else {
		responseMsg.Success = true
		responseMsg.Message = contents
	}
	return responseMsg
}

package rabbitmq

import (
	"encoding/json"
	"fmt"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/authstructs"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/eventingstructs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/loggingstructs"
	"github.com/MythicMeta/MythicContainer/translationstructs"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"github.com/MythicMeta/MythicContainer/webhookstructs"
	"os"
	"path/filepath"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	agentstructs.AllPayloadData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_WRITE_FILE,
		RabbitmqProcessingFunction: processC2RPCWriteFile,
	})
	authstructs.AllAuthData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_WRITE_FILE,
		RabbitmqProcessingFunction: processC2RPCWriteFile,
	})
	c2structs.AllC2Data.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_WRITE_FILE,
		RabbitmqProcessingFunction: processC2RPCWriteFile,
	})
	eventingstructs.AllEventingData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_WRITE_FILE,
		RabbitmqProcessingFunction: processC2RPCWriteFile,
	})
	loggingstructs.AllLoggingData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_WRITE_FILE,
		RabbitmqProcessingFunction: processC2RPCWriteFile,
	})
	translationstructs.AllTranslationData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_WRITE_FILE,
		RabbitmqProcessingFunction: processC2RPCWriteFile,
	})
	webhookstructs.AllWebhookData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_WRITE_FILE,
		RabbitmqProcessingFunction: processC2RPCWriteFile,
	})
}

func processC2RPCWriteFile(msg []byte) interface{} {
	input := sharedStructs.ContainerRPCWriteFileMessage{}
	responseMsg := sharedStructs.ContainerRPCWriteFileMessageResponse{
		Success: false,
		Error:   "Not implemented, not getting debug output",
	}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return ContainerRPCWriteFile(input)
	}
	return responseMsg
}

func ContainerRPCWriteFile(inputStruct sharedStructs.ContainerRPCWriteFileMessage) sharedStructs.ContainerRPCWriteFileMessageResponse {
	responseMsg := sharedStructs.ContainerRPCWriteFileMessageResponse{
		Success: false,
	}
	for _, containerName := range agentstructs.AllPayloadData.GetAllPayloadTypeNames() {
		if agentstructs.AllPayloadData.Get(containerName).GetPayloadDefinition().Name == inputStruct.ContainerName {
			return genericContainerWriteFile(inputStruct)
		}
	}
	for _, containerName := range c2structs.AllC2Data.GetAllNames() {
		if c2structs.AllC2Data.Get(containerName).GetC2Definition().Name == inputStruct.ContainerName {
			filePath, err := filepath.Abs(filepath.Join(c2structs.AllC2Data.Get(inputStruct.ContainerName).GetC2ServerFolderPath(), inputStruct.Filename))
			if err != nil {
				logging.LogError(err, "Failed to get absolute filepath for file to get")
				responseMsg.Error = fmt.Sprintf("Failed to locate file: %s\n", err.Error())
				return responseMsg
			}
			err = os.WriteFile(filePath, inputStruct.Contents, 0644)
			if err != nil {
				responseMsg.Error = fmt.Sprintf("Failed to open file: %s", err.Error())
				return responseMsg
			}
			responseMsg.Success = true
			responseMsg.Message = "Successfully wrote file"
			return responseMsg
		}
	}
	for _, containerName := range loggingstructs.AllLoggingData.GetAllNames() {
		if loggingstructs.AllLoggingData.Get(containerName).GetLoggingDefinition().Name == inputStruct.ContainerName {
			return genericContainerWriteFile(inputStruct)
		}
	}
	for _, containerName := range translationstructs.AllTranslationData.GetAllPayloadTypeNames() {
		if translationstructs.AllTranslationData.Get(containerName).GetPayloadDefinition().Name == inputStruct.ContainerName {
			return genericContainerWriteFile(inputStruct)
		}
	}
	for _, containerName := range webhookstructs.AllWebhookData.GetAllNames() {
		if webhookstructs.AllWebhookData.Get(containerName).GetWebhookDefinition().Name == inputStruct.ContainerName {
			return genericContainerWriteFile(inputStruct)
		}
	}
	for _, containerName := range eventingstructs.AllEventingData.GetAllNames() {
		if eventingstructs.AllEventingData.Get(containerName).GetEventingDefinition().Name == inputStruct.ContainerName {
			return genericContainerWriteFile(inputStruct)
		}
	}
	for _, containerName := range authstructs.AllAuthData.GetAllNames() {
		if authstructs.AllAuthData.Get(containerName).GetAuthDefinition().Name == inputStruct.ContainerName {
			return genericContainerWriteFile(inputStruct)
		}
	}
	responseMsg.Error = "failed to find appropriate container name"
	return responseMsg
}
func genericContainerWriteFile(inputStruct sharedStructs.ContainerRPCWriteFileMessage) sharedStructs.ContainerRPCWriteFileMessageResponse {
	responseMsg := sharedStructs.ContainerRPCWriteFileMessageResponse{
		Success: false,
	}
	filePath, err := filepath.Abs(filepath.Join(helpers.GetCwdFromExe(), inputStruct.Filename))
	if err != nil {
		logging.LogError(err, "Failed to get absolute filepath for file to get")
		responseMsg.Error = fmt.Sprintf("Failed to locate file: %s\n", err.Error())
		return responseMsg
	}
	err = os.WriteFile(filePath, inputStruct.Contents, 0644)
	if err != nil {
		responseMsg.Error = fmt.Sprintf("Failed to open file: %s", err.Error())
		return responseMsg
	}
	responseMsg.Success = true
	responseMsg.Message = "Successfully wrote file"
	return responseMsg
}

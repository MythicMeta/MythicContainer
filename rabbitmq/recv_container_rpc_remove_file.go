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
		RabbitmqRoutingKey:         CONTAINER_RPC_REMOVE_FILE,
		RabbitmqProcessingFunction: processContainerRPCRemoveFile,
	})
	authstructs.AllAuthData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_REMOVE_FILE,
		RabbitmqProcessingFunction: processContainerRPCRemoveFile,
	})
	c2structs.AllC2Data.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_REMOVE_FILE,
		RabbitmqProcessingFunction: processContainerRPCRemoveFile,
	})
	eventingstructs.AllEventingData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_REMOVE_FILE,
		RabbitmqProcessingFunction: processContainerRPCRemoveFile,
	})
	loggingstructs.AllLoggingData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_REMOVE_FILE,
		RabbitmqProcessingFunction: processContainerRPCRemoveFile,
	})
	translationstructs.AllTranslationData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_REMOVE_FILE,
		RabbitmqProcessingFunction: processContainerRPCRemoveFile,
	})
	webhookstructs.AllWebhookData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_REMOVE_FILE,
		RabbitmqProcessingFunction: processContainerRPCRemoveFile,
	})
}

func processContainerRPCRemoveFile(msg []byte) interface{} {
	input := sharedStructs.ContainerRPCRemoveFileMessage{}
	responseMsg := sharedStructs.ContainerRPCRemoveFileMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return ContainerRPCRemoveFile(input)
	}
	return responseMsg
}

func ContainerRPCRemoveFile(inputStruct sharedStructs.ContainerRPCRemoveFileMessage) sharedStructs.ContainerRPCRemoveFileMessageResponse {
	responseMsg := sharedStructs.ContainerRPCRemoveFileMessageResponse{
		Success: false,
	}
	for _, containerName := range agentstructs.AllPayloadData.GetAllPayloadTypeNames() {
		if agentstructs.AllPayloadData.Get(containerName).GetPayloadDefinition().Name == inputStruct.ContainerName {
			return genericContainerRemoveFile(inputStruct)
		}
	}
	for _, containerName := range c2structs.AllC2Data.GetAllNames() {
		if c2structs.AllC2Data.Get(containerName).GetC2Definition().Name == inputStruct.ContainerName {
			filePath, err := filepath.Abs(filepath.Join(c2structs.AllC2Data.Get(inputStruct.ContainerName).GetC2ServerFolderPath(), inputStruct.Filename))
			if err != nil {
				logging.LogError(err, "Failed to get absolute filepath for file to remove")
				responseMsg.Error = fmt.Sprintf("Failed to locate file: %s\n", err.Error())
			}
			err = os.Remove(filePath)
			if err != nil {
				responseMsg.Error = fmt.Sprintf("Failed to open file: %s", err.Error())
			}
			responseMsg.Success = true
			return responseMsg
		}
	}
	for _, containerName := range loggingstructs.AllLoggingData.GetAllNames() {
		if loggingstructs.AllLoggingData.Get(containerName).GetLoggingDefinition().Name == inputStruct.ContainerName {
			return genericContainerRemoveFile(inputStruct)
		}
	}
	for _, containerName := range translationstructs.AllTranslationData.GetAllPayloadTypeNames() {
		if translationstructs.AllTranslationData.Get(containerName).GetPayloadDefinition().Name == inputStruct.ContainerName {
			return genericContainerRemoveFile(inputStruct)
		}
	}
	for _, containerName := range webhookstructs.AllWebhookData.GetAllNames() {
		if webhookstructs.AllWebhookData.Get(containerName).GetWebhookDefinition().Name == inputStruct.ContainerName {
			return genericContainerRemoveFile(inputStruct)
		}
	}
	for _, containerName := range eventingstructs.AllEventingData.GetAllNames() {
		if eventingstructs.AllEventingData.Get(containerName).GetEventingDefinition().Name == inputStruct.ContainerName {
			return genericContainerRemoveFile(inputStruct)
		}
	}
	for _, containerName := range authstructs.AllAuthData.GetAllNames() {
		if authstructs.AllAuthData.Get(containerName).GetAuthDefinition().Name == inputStruct.ContainerName {
			return genericContainerRemoveFile(inputStruct)
		}
	}
	responseMsg.Error = "failed to find appropriate container name"
	return responseMsg
}
func genericContainerRemoveFile(inputStruct sharedStructs.ContainerRPCRemoveFileMessage) sharedStructs.ContainerRPCRemoveFileMessageResponse {
	responseMsg := sharedStructs.ContainerRPCRemoveFileMessageResponse{
		Success: false,
	}
	filePath, err := filepath.Abs(filepath.Join(helpers.GetCwdFromExe(), inputStruct.Filename))
	if err != nil {
		logging.LogError(err, "Failed to get absolute filepath for file to remove")
		responseMsg.Error = fmt.Sprintf("Failed to locate file: %s\n", err.Error())
	}
	err = os.Remove(filePath)
	if err != nil {
		responseMsg.Error = fmt.Sprintf("Failed to open file: %s", err.Error())
	}
	responseMsg.Success = true
	return responseMsg
}

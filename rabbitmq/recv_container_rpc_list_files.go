package rabbitmq

import (
	"encoding/json"
	"os"
	"path/filepath"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/authstructs"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/custombrowserstructs"
	"github.com/MythicMeta/MythicContainer/eventingstructs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/loggingstructs"
	"github.com/MythicMeta/MythicContainer/translationstructs"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"github.com/MythicMeta/MythicContainer/webhookstructs"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	agentstructs.AllPayloadData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_LIST_FILE,
		RabbitmqProcessingFunction: processContainerRPCListFile,
	})
	authstructs.AllAuthData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_LIST_FILE,
		RabbitmqProcessingFunction: processContainerRPCListFile,
	})
	c2structs.AllC2Data.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_LIST_FILE,
		RabbitmqProcessingFunction: processContainerRPCListFile,
	})
	eventingstructs.AllEventingData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_LIST_FILE,
		RabbitmqProcessingFunction: processContainerRPCListFile,
	})
	loggingstructs.AllLoggingData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_LIST_FILE,
		RabbitmqProcessingFunction: processContainerRPCListFile,
	})
	translationstructs.AllTranslationData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_LIST_FILE,
		RabbitmqProcessingFunction: processContainerRPCListFile,
	})
	webhookstructs.AllWebhookData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_LIST_FILE,
		RabbitmqProcessingFunction: processContainerRPCListFile,
	})
	custombrowserstructs.AllCustomBrowserData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONTAINER_RPC_LIST_FILE,
		RabbitmqProcessingFunction: processContainerRPCListFile,
	})
}

func processContainerRPCListFile(msg []byte) interface{} {
	input := sharedStructs.ContainerRPCListFileMessage{}
	responseMsg := sharedStructs.ContainerRPCListFileMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return ContainerRPCListFile(input)
	}
	return responseMsg
}

func ContainerRPCListFile(inputStruct sharedStructs.ContainerRPCListFileMessage) sharedStructs.ContainerRPCListFileMessageResponse {
	responseMsg := sharedStructs.ContainerRPCListFileMessageResponse{
		Success: false,
	}
	for _, containerName := range agentstructs.AllPayloadData.GetAllPayloadTypeNames() {
		if agentstructs.AllPayloadData.Get(containerName).GetPayloadDefinition().Name == inputStruct.ContainerName {
			return genericContainerListFiles(inputStruct)
		}
	}
	for _, containerName := range c2structs.AllC2Data.GetAllNames() {
		if c2structs.AllC2Data.Get(containerName).GetC2Definition().Name == inputStruct.ContainerName {
			path, err := filepath.Abs(c2structs.AllC2Data.Get(inputStruct.ContainerName).GetC2ServerFolderPath())
			if err != nil {
				logging.LogError(err, "Failed to get c2 server folder path")
				responseMsg.Error = err.Error()
				return responseMsg
			}
			entries, err := os.ReadDir(path)
			if err != nil {
				logging.LogError(err, "Failed to list out contents of server folder path")
				responseMsg.Error = err.Error()
				return responseMsg
			}
			logging.LogInfo("getting file list", "folder path", c2structs.AllC2Data.Get(inputStruct.ContainerName).GetC2ServerFolderPath())
			for _, entry := range entries {
				if !entry.IsDir() {
					responseMsg.Files = append(responseMsg.Files, entry.Name())
				}
			}
			responseMsg.Success = true
			return responseMsg
		}
	}
	for _, containerName := range loggingstructs.AllLoggingData.GetAllNames() {
		if loggingstructs.AllLoggingData.Get(containerName).GetLoggingDefinition().Name == inputStruct.ContainerName {
			return genericContainerListFiles(inputStruct)
		}
	}
	for _, containerName := range translationstructs.AllTranslationData.GetAllPayloadTypeNames() {
		if translationstructs.AllTranslationData.Get(containerName).GetPayloadDefinition().Name == inputStruct.ContainerName {
			return genericContainerListFiles(inputStruct)
		}
	}
	for _, containerName := range webhookstructs.AllWebhookData.GetAllNames() {
		if webhookstructs.AllWebhookData.Get(containerName).GetWebhookDefinition().Name == inputStruct.ContainerName {
			return genericContainerListFiles(inputStruct)
		}
	}
	for _, containerName := range eventingstructs.AllEventingData.GetAllNames() {
		if eventingstructs.AllEventingData.Get(containerName).GetEventingDefinition().Name == inputStruct.ContainerName {
			return genericContainerListFiles(inputStruct)
		}
	}
	for _, containerName := range authstructs.AllAuthData.GetAllNames() {
		if authstructs.AllAuthData.Get(containerName).GetAuthDefinition().Name == inputStruct.ContainerName {
			return genericContainerListFiles(inputStruct)
		}
	}
	for _, containerName := range custombrowserstructs.AllCustomBrowserData.GetAllNames() {
		if custombrowserstructs.AllCustomBrowserData.Get(containerName).GetCustomBrowserDefinition().Name == inputStruct.ContainerName {
			return genericContainerListFiles(inputStruct)
		}
	}
	responseMsg.Error = "failed to find appropriate container name"
	return responseMsg
}
func genericContainerListFiles(inputStruct sharedStructs.ContainerRPCListFileMessage) sharedStructs.ContainerRPCListFileMessageResponse {
	responseMsg := sharedStructs.ContainerRPCListFileMessageResponse{
		Success: false,
	}
	path, err := filepath.Abs(filepath.Join(helpers.GetCwdFromExe()))
	if err != nil {
		logging.LogError(err, "Failed to get c2 server folder path")
		responseMsg.Error = err.Error()
		return responseMsg
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		logging.LogError(err, "Failed to list out contents of server folder path")
		responseMsg.Error = err.Error()
		return responseMsg
	}
	logging.LogInfo("getting file list", "folder path", path)
	for _, entry := range entries {
		if !entry.IsDir() {
			responseMsg.Files = append(responseMsg.Files, entry.Name())
		}
	}
	responseMsg.Success = true
	return responseMsg
}

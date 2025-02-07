package rabbitmq

import (
	"encoding/json"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         PT_CHECK_IF_CALLBACKS_ALIVE,
		RabbitmqProcessingFunction: processPtCheckIfCallbacksAliveMessages,
	})
}

func processPtCheckIfCallbacksAliveMessages(msg []byte) interface{} {
	incomingMessage := agentstructs.PTCheckIfCallbacksAliveMessage{}
	response := agentstructs.PTCheckIfCallbacksAliveMessageResponse{}
	err := json.Unmarshal(msg, &incomingMessage)
	if err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Success = false
		response.Error = "Failed to unmarshal JSON message into structs"
		return response
	}
	response.Success = false
	checkIfCallbacksAliveFunc := agentstructs.AllPayloadData.Get(incomingMessage.ContainerName).GetCheckIfCallbacksAliveFunction()
	if checkIfCallbacksAliveFunc == nil {
		response.Success = true
	} else {
		response = checkIfCallbacksAliveFunc(incomingMessage)
	}
	return response
}

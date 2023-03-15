package rabbitmq

import (
	"encoding/json"
	"fmt"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddDirectMethod(agentstructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         PT_TASK_OPSEC_POST_CHECK,
		RabbitmqProcessingFunction: processPtTaskOPSECPostMessages,
	})
}

func processPtTaskOPSECPostMessages(msg []byte) {
	incomingMessage := agentstructs.PTTaskMessageAllData{}
	response := agentstructs.PTTaskOPSECPostTaskMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Error = "Failed to unmarshal JSON message into structs"
	} else {
		response.TaskID = incomingMessage.Task.ID
		for _, command := range agentstructs.AllPayloadData.Get(incomingMessage.PayloadType).GetCommands() {
			if command.Name == incomingMessage.Task.CommandName {
				if command.TaskFunctionOPSECPost != nil {
					if err := prepTaskArgs(command, &incomingMessage); err != nil {
						response.Error = err.Error()
						sendTaskOpsecPostResponse(response)
						return
					}
					response = command.TaskFunctionOPSECPost(&incomingMessage)
				} else {
					response.OpsecPostBlocked = false
					response.OpsecPostMessage = "Not Implemented"
					response.Success = true
				}
				sendTaskOpsecPostResponse(response)
				return
			}
		}
		response.Error = fmt.Sprintf("Failed to find command %s", incomingMessage.Task.CommandName)
		sendTaskOpsecPostResponse(response)
		return
	}
}

func sendTaskOpsecPostResponse(response agentstructs.PTTaskOPSECPostTaskMessageResponse) {
	if err := RabbitMQConnection.SendStructMessage(
		MYTHIC_EXCHANGE,
		PT_TASK_OPSEC_POST_CHECK_RESPONSE,
		"",
		response,
	); err != nil {
		logging.LogError(err, "Failed to send payload response back to Mythic")
	}
	return
}

package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         PT_TASK_OPSEC_PRE_CHECK,
		RabbitmqProcessingFunction: processPtTaskOPSECPreMessages,
	})
}

func processPtTaskOPSECPreMessages(msg []byte) {
	incomingMessage := agentstructs.PTTaskMessageAllData{}
	response := agentstructs.PTTTaskOPSECPreTaskMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Error = "Failed to unmarshal JSON message into structs"
	} else {
		response.TaskID = incomingMessage.Task.ID
		for _, command := range agentstructs.AllPayloadData.Get(incomingMessage.CommandPayloadType).GetCommands() {
			if command.Name == incomingMessage.Task.CommandName {
				if command.TaskFunctionOPSECPre != nil {
					if err := prepTaskArgs(command, &incomingMessage); err != nil {
						response.Error = err.Error()
						sendTaskOpsecPreResponse(response)
						return
					}
					response = command.TaskFunctionOPSECPre(&incomingMessage)
				} else {
					response.OpsecPreBlocked = false
					response.Success = true
					response.OpsecPreMessage = "Not Implemented"
				}
				response.TaskID = incomingMessage.Task.ID
				sendTaskOpsecPreResponse(response)
				return
			}
		}
		response.Error = fmt.Sprintf("Failed to find command %s", incomingMessage.Task.CommandName)
		sendTaskOpsecPreResponse(response)
		return
	}
}

func sendTaskOpsecPreResponse(response agentstructs.PTTTaskOPSECPreTaskMessageResponse) {
	for {
		err := RabbitMQConnection.SendStructMessage(
			MYTHIC_EXCHANGE,
			PT_TASK_OPSEC_PRE_CHECK_RESPONSE,
			"",
			response,
			false,
		)
		if err != nil {
			logging.LogError(err, "Failed to send payload response back to Mythic")
			continue
		}
		return
	}

}

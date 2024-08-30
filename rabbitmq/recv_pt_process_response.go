package rabbitmq

import (
	"encoding/json"
	"fmt"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         PT_TASK_PROCESS_RESPONSE,
		RabbitmqProcessingFunction: processPtProcessResponseMessages,
	})
}

func processPtProcessResponseMessages(msg []byte) {
	incomingMessage := agentstructs.PtTaskProcessResponseMessage{}
	response := agentstructs.PTTaskProcessResponseMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Error = "Failed to unmarshal JSON message into structs"
		sendTaskProcessResponseResponse(response)
		return
	} else {
		for _, command := range agentstructs.AllPayloadData.Get(incomingMessage.TaskData.CommandPayloadType).GetCommands() {
			if command.Name == incomingMessage.TaskData.Task.CommandName {
				if err := prepTaskArgs(command, incomingMessage.TaskData); err != nil {
					response.Error = err.Error()
					sendTaskProcessResponseResponse(response)
					return
				}
				if command.TaskFunctionProcessResponse != nil {
					response = command.TaskFunctionProcessResponse(incomingMessage)
				} else {
					response.Error = fmt.Sprintf("Failed to find process response function for command %s", incomingMessage.TaskData.Task.CommandName)
				}
				response.TaskID = incomingMessage.TaskData.Task.ID
				sendTaskProcessResponseResponse(response)
				return
			}
		}
		response.Error = fmt.Sprintf("Failed to find command %s", incomingMessage.TaskData.Task.CommandName)
		sendTaskProcessResponseResponse(response)
		return
	}
}

func sendTaskProcessResponseResponse(response agentstructs.PTTaskProcessResponseMessageResponse) {
	if err := RabbitMQConnection.SendStructMessage(
		MYTHIC_EXCHANGE,
		PT_TASK_PROCESS_RESPONSE_RESPONSE,
		"",
		response,
		false,
	); err != nil {
		logging.LogError(err, "Failed to send payload response back to Mythic")
	}
}

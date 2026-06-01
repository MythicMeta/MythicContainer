package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         PT_TASK_COMPLETION_FUNCTION,
		RabbitmqProcessingFunction: processPtTaskCompletionMessages,
	})
}

func processPtTaskCompletionMessages(ctx context.Context, msg []byte) {
	incomingMessage := agentstructs.PTTaskCompletionFunctionMessage{}
	response := agentstructs.PTTaskCompletionFunctionMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Error = "Failed to unmarshal JSON message into structs"
		sendTaskCompletionResponse(ctx, response)
		return
	} else {
		for _, command := range agentstructs.AllPayloadData.Get(incomingMessage.TaskData.CommandPayloadType).GetCommands() {
			if command.Name == incomingMessage.TaskData.Task.CommandName {
				if err := prepTaskArgs(ctx, command, incomingMessage.TaskData); err != nil {
					response.Error = err.Error()
					sendTaskCompletionResponse(ctx, response)
					return
				}
				if command.TaskCompletionFunctions != nil {
					found := false
					for funcName, funcDef := range command.TaskCompletionFunctions {
						if funcName == incomingMessage.CompletionFunctionName {
							found = true
							response = funcDef(ctx, incomingMessage.TaskData, incomingMessage.SubtaskData, incomingMessage.SubtaskGroup)
						}
					}
					if !found {
						response.Error = fmt.Sprintf("Failed to find completion function: %s", incomingMessage.CompletionFunctionName)
					}
				} else {
					response.Error = fmt.Sprintf("Failed to find completion function: %s", incomingMessage.CompletionFunctionName)
				}
				if incomingMessage.SubtaskData != nil {
					response.TaskID = incomingMessage.SubtaskData.Task.ID
					response.ParentTaskId = incomingMessage.TaskData.Task.ID
				} else {
					response.TaskID = incomingMessage.TaskData.Task.ID
				}
				sendTaskCompletionResponse(ctx, response)
				return
			}
		}
		response.Error = fmt.Sprintf("Failed to find command %s", incomingMessage.TaskData.Task.CommandName)
		sendTaskCompletionResponse(ctx, response)
		return
	}
}

func sendTaskCompletionResponse(ctx context.Context, response agentstructs.PTTaskCompletionFunctionMessageResponse) {
	for {
		err := RabbitMQConnection.SendStructMessageWithContext(ctx,
			MYTHIC_EXCHANGE,
			PT_TASK_COMPLETION_FUNCTION_RESPONSE,
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

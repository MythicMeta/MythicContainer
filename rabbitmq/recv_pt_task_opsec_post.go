package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         PT_TASK_OPSEC_POST_CHECK,
		RabbitmqProcessingFunction: processPtTaskOPSECPostMessages,
	})
}

func processPtTaskOPSECPostMessages(ctx context.Context, msg []byte) {
	incomingMessage := agentstructs.PTTaskMessageAllData{}
	response := agentstructs.PTTaskOPSECPostTaskMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Error = "Failed to unmarshal JSON message into structs"
	} else {
		response.TaskID = incomingMessage.Task.ID
		for _, command := range agentstructs.AllPayloadData.Get(incomingMessage.CommandPayloadType).GetCommands() {
			if command.Name == incomingMessage.Task.CommandName {
				if command.TaskFunctionOPSECPost != nil {
					if err := prepTaskArgs(ctx, command, &incomingMessage); err != nil {
						response.Error = err.Error()
						sendTaskOpsecPostResponse(ctx, response)
						return
					}
					response = command.TaskFunctionOPSECPost(ctx, &incomingMessage)
				} else {
					response.OpsecPostBlocked = false
					response.OpsecPostMessage = "Not Implemented"
					response.Success = true
				}
				sendTaskOpsecPostResponse(ctx, response)
				return
			}
		}
		response.Error = fmt.Sprintf("Failed to find command %s", incomingMessage.Task.CommandName)
		sendTaskOpsecPostResponse(ctx, response)
		return
	}
}

func sendTaskOpsecPostResponse(ctx context.Context, response agentstructs.PTTaskOPSECPostTaskMessageResponse) {
	for {
		err := RabbitMQConnection.SendStructMessageWithContext(ctx,
			MYTHIC_EXCHANGE,
			PT_TASK_OPSEC_POST_CHECK_RESPONSE,
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

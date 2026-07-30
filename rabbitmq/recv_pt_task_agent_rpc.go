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
		RabbitmqRoutingKey:         PT_TASK_AGENT_RPC,
		RabbitmqProcessingFunction: processPtTaskAgentRPCMessages,
	})
}

func processPtTaskAgentRPCMessages(ctx context.Context, msg []byte) {
	request := agentstructs.PTTaskAgentRPCMessage{}
	if err := json.Unmarshal(msg, &request); err != nil {
		logging.LogError(err, "Failed to unmarshal agent RPC request")
		return
	}
	response := invokePTTaskAgentRPC(ctx, request)
	response.CallbackID = request.TaskData.Callback.ID
	response.AgentTaskID = request.TaskData.Task.AgentTaskID
	if response.Status == "" {
		response.Status = "success"
	}
	sendPTTaskAgentRPCResponse(ctx, response)
}

func invokePTTaskAgentRPC(ctx context.Context, incomingMessage agentstructs.PTTaskAgentRPCMessage) (response agentstructs.PTTaskAgentRPCMessageResponse) {
	response.Status = "error"
	if incomingMessage.TaskData == nil {
		response.Output = "agent RPC request missing task"
		return response
	}
	for _, command := range agentstructs.AllPayloadData.Get(incomingMessage.TaskData.CommandPayloadType).GetCommands() {
		if command.Name == incomingMessage.TaskData.Task.CommandName {
			if err := prepTaskArgs(ctx, command, incomingMessage.TaskData); err != nil {
				response.Output = err.Error()
				return response
			}
			if command.TaskFunctionProcessResponse != nil {
				response = command.AgentRPCFunction(ctx, incomingMessage.TaskData, incomingMessage.Name, incomingMessage.Arguments)
			} else {
				response.Output = fmt.Sprintf("Failed to find process response function for command %s", incomingMessage.TaskData.Task.CommandName)
			}
			return response
		}
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			response = agentstructs.PTTaskAgentRPCMessageResponse{
				Status: "error",
				Output: fmt.Sprintf("agent RPC function panicked: %v", recovered),
			}
		}
	}()
	return response
}

func sendPTTaskAgentRPCResponse(ctx context.Context, response agentstructs.PTTaskAgentRPCMessageResponse) {
	for {
		err := RabbitMQConnection.SendStructMessageWithContext(
			ctx,
			MYTHIC_EXCHANGE,
			PT_TASK_AGENT_RPC_RESPONSE,
			"",
			response,
			false,
		)
		if err != nil {
			logging.LogError(err, "Failed to send agent RPC response back to Mythic")
			if ctx.Err() != nil {
				return
			}
			continue
		}
		return
	}
}

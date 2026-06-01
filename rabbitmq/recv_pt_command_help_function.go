package rabbitmq

import (
	"context"
	"encoding/json"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         PT_RPC_COMMAND_HELP_FUNCTION,
		RabbitmqProcessingFunction: processPtRPCCommandHelpFunctionMessages,
	})
}

func processPtRPCCommandHelpFunctionMessages(ctx context.Context, msg []byte) interface{} {
	incomingMessage := agentstructs.PTRPCCommandHelpFunctionMessage{}
	response := agentstructs.PTRPCCommandHelpFunctionMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Error = "Failed to unmarshal JSON message into structs"
		return response
	} else {
		pt := agentstructs.AllPayloadData.Get(incomingMessage.PayloadType).GetPayloadDefinition()
		if pt.CommandHelpFunction != nil {
			resp := pt.CommandHelpFunction(ctx, incomingMessage)
			return resp
		}
		response.Error = "Function is null"
		return response
	}
}

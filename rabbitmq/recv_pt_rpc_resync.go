package rabbitmq

import (
	"encoding/json"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	agentstructs.AllPayloadData.Get("").AddRPCMethod(agentstructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         PT_RPC_RESYNC_ROUTING_KEY,
		RabbitmqProcessingFunction: processPTRPCReSync,
	})
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processPTRPCReSync(msg []byte) interface{} {
	input := agentstructs.PTRPCReSyncMessage{}
	responseMsg := agentstructs.PTRPCReSyncMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return PTRPCReSync(input)
	}
	return responseMsg
}

func PTRPCReSync(input agentstructs.PTRPCReSyncMessage) agentstructs.PTRPCReSyncMessageResponse {
	response := agentstructs.PTRPCReSyncMessageResponse{
		Success: true,
		Error:   "",
	}
	SyncPayloadData(&input.Name)
	return response
}

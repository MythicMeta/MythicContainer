package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/translationstructs"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	translationstructs.AllTranslationData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         TR_RPC_RESYNC_ROUTING_KEY,
		RabbitmqProcessingFunction: processTrRPCReSync,
	})
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processTrRPCReSync(msg []byte) interface{} {
	input := translationstructs.TRRPCReSyncMessage{}
	responseMsg := translationstructs.TRRPCReSyncMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return TrRPCReSync(input)
	}
	return responseMsg
}

func TrRPCReSync(input translationstructs.TRRPCReSyncMessage) translationstructs.TRRPCReSyncMessageResponse {
	response := translationstructs.TRRPCReSyncMessageResponse{
		Success: true,
		Error:   "",
	}
	SyncTranslationData(&input.Name)
	return response
}

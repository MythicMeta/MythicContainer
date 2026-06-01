package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/translationstructs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	/*
		translationstructs.AllTranslationData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
			RabbitmqRoutingKey:         TR_RPC_ENCRYPT_BYTES,
			RabbitmqProcessingFunction: processTrRPCEncryptBytes,
		})

	*/
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processTrRPCEncryptBytes(ctx context.Context, msg []byte) interface{} {
	input := translationstructs.TrEncryptBytesMessage{}
	responseMsg := translationstructs.TrEncryptBytesMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return TrRPCEncryptBytes(ctx, input)
	}
	return responseMsg
}

func TrRPCEncryptBytes(ctx context.Context, input translationstructs.TrEncryptBytesMessage) translationstructs.TrEncryptBytesMessageResponse {
	response := translationstructs.TrEncryptBytesMessageResponse{
		Success: false,
		Error:   "No Translation function defined",
	}
	if translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().EncryptBytes != nil {
		response = translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().EncryptBytes(ctx, input)
	}
	return response
}

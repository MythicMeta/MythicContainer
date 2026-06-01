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
			RabbitmqRoutingKey:         TR_RPC_CONVERT_FROM_MYTHIC_C2_FORMAT,
			RabbitmqProcessingFunction: processTrRPCMythicToCustomFormat,
		})

	*/
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processTrRPCMythicToCustomFormat(ctx context.Context, msg []byte) interface{} {
	input := translationstructs.TrMythicC2ToCustomMessageFormatMessage{}
	responseMsg := translationstructs.TrMythicC2ToCustomMessageFormatMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return TrRPCMythicToCustomFormat(ctx, input)
	}
	return responseMsg
}

func TrRPCMythicToCustomFormat(ctx context.Context, input translationstructs.TrMythicC2ToCustomMessageFormatMessage) translationstructs.TrMythicC2ToCustomMessageFormatMessageResponse {
	response := translationstructs.TrMythicC2ToCustomMessageFormatMessageResponse{
		Success: false,
		Error:   "No Translation function defined",
	}
	if translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().TranslateMythicToCustomFormat != nil {
		response = translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().TranslateMythicToCustomFormat(ctx, input)
	}
	return response
}

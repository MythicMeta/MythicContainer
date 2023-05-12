package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/translationstructs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	/*
		translationstructs.AllTranslationData.Get("").AddRPCMethod(translationstructs.RabbitmqRPCMethod{
			RabbitmqRoutingKey:         TR_RPC_CONVERT_TO_MYTHIC_C2_FORMAT,
			RabbitmqProcessingFunction: processTrRPCCustomToMythicFormat,
		})

	*/
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processTrRPCCustomToMythicFormat(msg []byte) interface{} {
	input := translationstructs.TrCustomMessageToMythicC2FormatMessage{}
	responseMsg := translationstructs.TrCustomMessageToMythicC2FormatMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return TrRPCCustomToMythicFormat(input)
	}
	return responseMsg
}

func TrRPCCustomToMythicFormat(input translationstructs.TrCustomMessageToMythicC2FormatMessage) translationstructs.TrCustomMessageToMythicC2FormatMessageResponse {
	response := translationstructs.TrCustomMessageToMythicC2FormatMessageResponse{
		Success: false,
		Error:   "No Translation function defined",
	}
	if translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().TranslateCustomToMythicFormat != nil {
		response = translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().TranslateCustomToMythicFormat(input)
	}
	return response
}

package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/translationstructs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	translationstructs.AllTranslationData.Get("").AddRPCMethod(translationstructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         TR_RPC_GENERATE_KEYS,
		RabbitmqProcessingFunction: processTrRPCGenerateEncryptionKeys,
	})
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processTrRPCGenerateEncryptionKeys(msg []byte) interface{} {
	input := translationstructs.TrGenerateEncryptionKeysMessage{}
	responseMsg := translationstructs.TrGenerateEncryptionKeysMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return TrRPCGenerateEncryptionKeys(input)
	}
	return responseMsg
}

func TrRPCGenerateEncryptionKeys(input translationstructs.TrGenerateEncryptionKeysMessage) translationstructs.TrGenerateEncryptionKeysMessageResponse {
	response := translationstructs.TrGenerateEncryptionKeysMessageResponse{
		Success: false,
		Error:   "No Translation function defined",
	}
	//logging.LogDebug("asked to generate keys", "req", input)
	if translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().GenerateEncryptionKeys != nil {
		response = translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().GenerateEncryptionKeys(input)
	}
	return response
}

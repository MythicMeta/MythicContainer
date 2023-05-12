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
			RabbitmqRoutingKey:         TR_RPC_DECRYPT_BYTES,
			RabbitmqProcessingFunction: processTrRPCDecryptBytes,
		})

	*/
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processTrRPCDecryptBytes(msg []byte) interface{} {
	input := translationstructs.TrDecryptBytesMessage{}
	responseMsg := translationstructs.TrDecryptBytesMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		return TrRPCDecryptBytes(input)
	}
	return responseMsg
}

func TrRPCDecryptBytes(input translationstructs.TrDecryptBytesMessage) translationstructs.TrDecryptBytesMessageResponse {
	response := translationstructs.TrDecryptBytesMessageResponse{
		Success: false,
		Error:   "No Translation function defined",
	}
	if translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().DecryptBytes != nil {
		response = translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().DecryptBytes(input)
	}
	return response
}

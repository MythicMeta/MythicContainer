package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/authstructs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"slices"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	authstructs.AllAuthData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         AUTH_RPC_GET_IDP_METADATA,
		RabbitmqProcessingFunction: processAuthRPCGetIDPMetadata,
	})
}

func processAuthRPCGetIDPMetadata(msg []byte) interface{} {
	input := authstructs.GetIDPMetadataMessage{}
	responseMsg := authstructs.GetIDPMetadataMessageResponse{
		Success: false,
		Error:   "Not implemented, not getting metadata",
	}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return AuthRPCGetIDPMetadata(input)
	}
	return responseMsg
}

func AuthRPCGetIDPMetadata(input authstructs.GetIDPMetadataMessage) authstructs.GetIDPMetadataMessageResponse {
	responseMsg := authstructs.GetIDPMetadataMessageResponse{
		Success: false,
		Error:   "Failed to find right container with a non-nil function",
	}
	for _, eventing := range authstructs.AllAuthData.GetAllNames() {
		if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().Name == input.ContainerName {
			if slices.Contains(authstructs.AllAuthData.Get(eventing).GetAuthDefinition().IDPServices, input.IDPName) {
				if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().GetIDPMetadata != nil {
					response := authstructs.AllAuthData.Get(eventing).GetAuthDefinition().GetIDPMetadata(input)
					return response
				}
			}
		}
	}
	return responseMsg
}

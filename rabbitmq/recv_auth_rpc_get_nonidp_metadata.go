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
		RabbitmqRoutingKey:         AUTH_RPC_GET_NONIDP_METADATA,
		RabbitmqProcessingFunction: processAuthRPCGetNonIDPMetadata,
	})
}

func processAuthRPCGetNonIDPMetadata(msg []byte) interface{} {
	input := authstructs.GetNonIDPMetadataMessage{}
	responseMsg := authstructs.GetNonIDPMetadataMessageResponse{
		Success: false,
		Error:   "Not implemented, not getting metadata",
	}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return AuthRPCGetNonIDPMetadata(input)
	}
	return responseMsg
}

func AuthRPCGetNonIDPMetadata(input authstructs.GetNonIDPMetadataMessage) authstructs.GetNonIDPMetadataMessageResponse {
	responseMsg := authstructs.GetNonIDPMetadataMessageResponse{
		Success: false,
		Error:   "Failed to find right container with a non-nil function",
	}
	for _, eventing := range authstructs.AllAuthData.GetAllNames() {
		if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().Name == input.ContainerName {
			if slices.Contains(authstructs.AllAuthData.Get(eventing).GetAuthDefinition().NonIDPServices, input.NonIDPName) {
				if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().GetNonIDPMetadata != nil {
					response := authstructs.AllAuthData.Get(eventing).GetAuthDefinition().GetNonIDPMetadata(input)
					return response
				}
			}
		}
	}
	return responseMsg
}

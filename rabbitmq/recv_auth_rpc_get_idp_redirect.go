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
		RabbitmqRoutingKey:         AUTH_RPC_GET_IDP_REDIRECT,
		RabbitmqProcessingFunction: processAuthRPCGetIDPRedirect,
	})
}

func processAuthRPCGetIDPRedirect(msg []byte) interface{} {
	input := authstructs.GetIDPRedirectMessage{}
	responseMsg := authstructs.GetIDPRedirectMessageResponse{
		Success: false,
		Error:   "Not implemented, not getting debug output",
	}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return AuthRPCGetIDPRedirect(input)
	}
	return responseMsg
}

func AuthRPCGetIDPRedirect(input authstructs.GetIDPRedirectMessage) authstructs.GetIDPRedirectMessageResponse {
	responseMsg := authstructs.GetIDPRedirectMessageResponse{
		Success: false,
		Error:   "Failed to find right container with a non-nil function",
	}
	for _, eventing := range authstructs.AllAuthData.GetAllNames() {
		if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().Name == input.ContainerName {
			if slices.Contains(authstructs.AllAuthData.Get(eventing).GetAuthDefinition().IDPServices, input.IDPName) {
				if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().GetIDPRedirect != nil {
					return authstructs.AllAuthData.Get(eventing).GetAuthDefinition().GetIDPRedirect(input)
				}
			}

		}
	}
	return responseMsg
}

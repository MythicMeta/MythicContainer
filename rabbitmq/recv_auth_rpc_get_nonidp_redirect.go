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
		RabbitmqRoutingKey:         AUTH_RPC_GET_NONIDP_REDIRECT,
		RabbitmqProcessingFunction: processAuthRPCGetNonIDPRedirect,
	})
}

func processAuthRPCGetNonIDPRedirect(msg []byte) interface{} {
	input := authstructs.GetNonIDPRedirectMessage{}
	responseMsg := authstructs.GetNonIDPRedirectMessageResponse{
		Success: false,
		Error:   "Not implemented, not getting debug output",
	}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return AuthRPCGetNonIDPRedirect(input)
	}
	return responseMsg
}

func AuthRPCGetNonIDPRedirect(input authstructs.GetNonIDPRedirectMessage) authstructs.GetNonIDPRedirectMessageResponse {
	responseMsg := authstructs.GetNonIDPRedirectMessageResponse{
		Success: false,
		Error:   "Failed to find right container with a non-nil function",
	}
	for _, eventing := range authstructs.AllAuthData.GetAllNames() {
		if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().Name == input.ContainerName {
			if slices.Contains(authstructs.AllAuthData.Get(eventing).GetAuthDefinition().NonIDPServices, input.NonIDPName) {
				if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().GetNonIDPRedirect != nil {
					return authstructs.AllAuthData.Get(eventing).GetAuthDefinition().GetNonIDPRedirect(input)
				}
			}

		}
	}
	return responseMsg
}

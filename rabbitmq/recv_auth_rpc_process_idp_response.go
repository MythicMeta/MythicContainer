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
		RabbitmqRoutingKey:         AUTH_RPC_PROCESS_IDP_RESPONSE,
		RabbitmqProcessingFunction: processAuthRPCProcessIDPResponse,
	})
}

func processAuthRPCProcessIDPResponse(msg []byte) interface{} {
	input := authstructs.ProcessIDPResponseMessage{}
	responseMsg := authstructs.ProcessIDPResponseMessageResponse{
		SuccessfulAuthentication: false,
		Error:                    "Not implemented, not authing",
	}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.SuccessfulAuthentication = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return AuthRPCProcessIDPResponse(input)
	}
	return responseMsg
}

func AuthRPCProcessIDPResponse(input authstructs.ProcessIDPResponseMessage) authstructs.ProcessIDPResponseMessageResponse {
	responseMsg := authstructs.ProcessIDPResponseMessageResponse{
		SuccessfulAuthentication: false,
		Error:                    "Failed to find right container with a non-nil function",
	}
	for _, eventing := range authstructs.AllAuthData.GetAllNames() {
		if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().Name == input.ContainerName {
			if slices.Contains(authstructs.AllAuthData.Get(eventing).GetAuthDefinition().IDPServices, input.IDPName) {
				if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().ProcessIDPResponse != nil {
					return authstructs.AllAuthData.Get(eventing).GetAuthDefinition().ProcessIDPResponse(input)
				}
			}

		}
	}
	return responseMsg
}

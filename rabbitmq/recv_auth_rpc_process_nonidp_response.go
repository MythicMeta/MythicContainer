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
		RabbitmqRoutingKey:         AUTH_RPC_PROCESS_NONIDP_RESPONSE,
		RabbitmqProcessingFunction: processAuthRPCProcessNonIDPResponse,
	})
}

func processAuthRPCProcessNonIDPResponse(msg []byte) interface{} {
	input := authstructs.ProcessNonIDPResponseMessage{}
	responseMsg := authstructs.ProcessNonIDPResponseMessageResponse{
		SuccessfulAuthentication: false,
		Error:                    "Not implemented, not authing",
	}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.SuccessfulAuthentication = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return AuthRPCProcessNonIDPResponse(input)
	}
	return responseMsg
}

func AuthRPCProcessNonIDPResponse(input authstructs.ProcessNonIDPResponseMessage) authstructs.ProcessNonIDPResponseMessageResponse {
	responseMsg := authstructs.ProcessNonIDPResponseMessageResponse{
		SuccessfulAuthentication: false,
		Error:                    "Failed to find right container with a non-nil function",
	}
	for _, eventing := range authstructs.AllAuthData.GetAllNames() {
		if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().Name == input.ContainerName {
			if slices.Contains(authstructs.AllAuthData.Get(eventing).GetAuthDefinition().NonIDPServices, input.NonIDPName) {
				if authstructs.AllAuthData.Get(eventing).GetAuthDefinition().ProcessNonIDPResponse != nil {
					return authstructs.AllAuthData.Get(eventing).GetAuthDefinition().ProcessNonIDPResponse(input)
				}
			}

		}
	}
	return responseMsg
}

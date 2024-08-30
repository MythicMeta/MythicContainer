package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/authstructs"
	"github.com/MythicMeta/MythicContainer/eventingstructs"
	"github.com/MythicMeta/MythicContainer/loggingstructs"
	"github.com/MythicMeta/MythicContainer/translationstructs"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"github.com/MythicMeta/MythicContainer/webhookstructs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	loggingstructs.AllLoggingData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONSUMING_CONTAINER_RESYNC_ROUTING_KEY,
		RabbitmqProcessingFunction: processConsumingServiceRPCReSync,
	})
	webhookstructs.AllWebhookData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONSUMING_CONTAINER_RESYNC_ROUTING_KEY,
		RabbitmqProcessingFunction: processConsumingServiceRPCReSync,
	})
	eventingstructs.AllEventingData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONSUMING_CONTAINER_RESYNC_ROUTING_KEY,
		RabbitmqProcessingFunction: processConsumingServiceRPCReSync,
	})
	authstructs.AllAuthData.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         CONSUMING_CONTAINER_RESYNC_ROUTING_KEY,
		RabbitmqProcessingFunction: processConsumingServiceRPCReSync,
	})
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processConsumingServiceRPCReSync(msg []byte) interface{} {
	input := translationstructs.TRRPCReSyncMessage{}
	responseMsg := translationstructs.TRRPCReSyncMessageResponse{Success: true}
	err := json.Unmarshal(msg, &input)
	if err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	}
	return ConsumingContainerReSync(input)
}

func ConsumingContainerReSync(input translationstructs.TRRPCReSyncMessage) translationstructs.TRRPCReSyncMessageResponse {
	response := translationstructs.TRRPCReSyncMessageResponse{
		Success: true,
		Error:   "",
	}
	if loggingstructs.AllLoggingData.Get(input.Name).GetLoggingDefinition().Name == input.Name {
		SyncConsumingContainerData(input.Name, "logging")
	}
	if webhookstructs.AllWebhookData.Get(input.Name).GetWebhookDefinition().Name == input.Name {
		SyncConsumingContainerData(input.Name, "webhook")
	}
	if eventingstructs.AllEventingData.Get(input.Name).GetEventingDefinition().Name == input.Name {
		SyncConsumingContainerData(input.Name, "eventing")
	}
	if authstructs.AllAuthData.Get(input.Name).GetAuthDefinition().Name == input.Name {
		SyncConsumingContainerData(input.Name, "auth")
	}
	return response
}

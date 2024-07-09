package rabbitmq

import (
	"encoding/json"
	"errors"
	"github.com/MythicMeta/MythicContainer/authstructs"
	"github.com/MythicMeta/MythicContainer/eventingstructs"
	"github.com/MythicMeta/MythicContainer/loggingstructs"
	"github.com/MythicMeta/MythicContainer/webhookstructs"
	"time"

	"github.com/MythicMeta/MythicContainer/logging"
)

type ConsumingContainerSyncResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SyncConsumingContainerData(consumingContainerName string, consumingType string) {
	if consumingContainerName == "" {
		return
	}
	logging.LogInfo("Syncing consuming container", "name", consumingContainerName)
	response := ConsumingContainerSyncResponse{}
	var description string
	var subscriptions []string
	switch consumingType {
	case "logging":
		def := loggingstructs.AllLoggingData.Get(consumingContainerName).GetLoggingDefinition()
		if def.Name == "" {
			logging.LogError(nil, "Failed to find logging info to sync")
			return
		}
		description = def.Description
		subscriptions = def.Subscriptions
	case "webhook":
		def := webhookstructs.AllWebhookData.Get(consumingContainerName).GetWebhookDefinition()
		if def.Name == "" {
			logging.LogError(nil, "Failed to find webhook info to sync")
			return
		}
		description = def.Description
		subscriptions = def.Subscriptions
	case "eventing":
		def := eventingstructs.AllEventingData.Get(consumingContainerName).GetEventingDefinition()
		if def.Name == "" {
			logging.LogError(nil, "Failed to find eventing info to sync")
			return
		}
		description = def.Description
		subscriptions = def.Subscriptions
	case "auth":
		def := authstructs.AllAuthData.Get(consumingContainerName).GetAuthDefinition()
		if def.Name == "" {
			logging.LogError(nil, "Failed to find eventing info to sync")
			return
		}
		description = def.Description
		subscriptions = def.Subscriptions
	default:
	}
	syncMessage := map[string]interface{}{
		"consuming_container": map[string]interface{}{
			"name":          consumingContainerName,
			"type":          consumingType,
			"description":   description,
			"subscriptions": subscriptions,
		},
		"container_version": containerVersion,
	}
	for {
		syncMessageJson, err := json.Marshal(syncMessage)
		if err != nil {
			logging.LogError(err, "Failed to serialize consuming container message to json, %s", err.Error())
			time.Sleep(RETRY_CONNECT_DELAY)
		}
		resp, err := RabbitMQConnection.SendRPCMessage(MYTHIC_EXCHANGE, CONSUMING_CONTAINER_SYNC_ROUTING_KEY, syncMessageJson, true)
		if err != nil {
			logging.LogError(err, "Failed to send consuming container to Mythic")
			time.Sleep(RETRY_CONNECT_DELAY)
		}
		err = json.Unmarshal(resp, &response)
		if err != nil {
			logging.LogError(err, "Failed to marshal sync response back to struct")
			time.Sleep(RETRY_CONNECT_DELAY)
		}
		if !response.Success {
			logging.LogError(errors.New(response.Error), "waiting and trying again...")
			time.Sleep(RETRY_CONNECT_DELAY)
		}
		logging.LogInfo("Successfully synced consuming container!", "name", consumingContainerName,
			"type", consumingType)
		break
	}
}

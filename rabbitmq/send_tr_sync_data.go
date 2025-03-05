package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/translationstructs"
	"time"

	"github.com/MythicMeta/MythicContainer/logging"
)

func SyncTranslationData(translationName *string) {
	// now make our payloadtype info that we're going to sync
	for _, pt := range translationstructs.AllTranslationData.GetAllPayloadTypeNames() {
		if translationName == nil || *translationName == pt {
			logging.LogInfo("Syncing translation container", "name", pt)
			syncMessage := translationstructs.TrSyncMessage{}
			response := translationstructs.TrSyncMessageResponse{}
			syncMessage.Name = translationstructs.AllTranslationData.Get(pt).GetPayloadName()
			syncMessage.Author = translationstructs.AllTranslationData.Get(pt).GetAuthor()
			syncMessage.Description = translationstructs.AllTranslationData.Get(pt).GetDescription()
			translationstructs.AllTranslationData.Get(pt).AddContainerVersion(containerVersion)
			syncMessage.ContainerVersion = translationstructs.AllTranslationData.Get(pt).GetContainerVersion()
			//logging.LogDebug("syncing over tr", "tr info", syncMessage)
			for {
				syncMessageJson, err := json.Marshal(syncMessage)
				if err != nil {
					logging.LogError(err, "Failed to serialize translation service sync message to json, %s", err.Error())
					time.Sleep(RETRY_CONNECT_DELAY)
					continue
				}
				resp, err := RabbitMQConnection.SendRPCMessage(MYTHIC_EXCHANGE, TR_SYNC_ROUTING_KEY, syncMessageJson, true)
				if err != nil {
					logging.LogError(err, "Failed to send translation service to Mythic")
					time.Sleep(RETRY_CONNECT_DELAY)
					continue
				}
				err = json.Unmarshal(resp, &response)
				if err != nil {
					logging.LogError(err, "Failed to marshal sync response back to struct")
					time.Sleep(RETRY_CONNECT_DELAY)
					continue
				}
				if !response.Success {
					logging.LogError(nil, response.Error)
					time.Sleep(RETRY_CONNECT_DELAY)
					continue
				}
				logging.LogInfo("Successfully synced translation service!", "name", pt)
				break
			}
		}
	}
}

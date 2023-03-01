package rabbitmq

import (
	"encoding/json"
	c2structs "github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"time"
)

func SyncAllC2Data(resyncName *string) {
	// now make our c2 info that we're going to sync
	for _, c2 := range c2structs.AllC2Data.GetAllNames() {
		if resyncName == nil || c2 == *resyncName {
			logging.LogInfo("Syncing C2 profile", "name", c2)
			syncMessage := c2structs.C2SyncMessage{}
			response := c2structs.C2SyncMessageResponse{}
			syncMessage.Profile = c2structs.AllC2Data.Get(c2).GetC2Definition()
			syncMessage.Parameters = c2structs.AllC2Data.Get(c2).GetParameters()
			c2structs.AllC2Data.Get(c2).AddContainerVersion(containerVersion)
			syncMessage.ContainerVersion = c2structs.AllC2Data.Get(c2).GetContainerVersion()
			for {
				if syncMessageJson, err := json.Marshal(syncMessage); err != nil {
					logging.LogError(err, "Failed to serialize c2 sync message to json")
					time.Sleep(RETRY_CONNECT_DELAY)
				} else if resp, err := RabbitMQConnection.SendRPCMessage(MYTHIC_EXCHANGE, C2_SYNC_ROUTING_KEY, syncMessageJson, true); err != nil {
					logging.LogError(err, "Failed to send c2 profile to Mythic")
					time.Sleep(RETRY_CONNECT_DELAY)
				} else if err := json.Unmarshal(resp, &response); err != nil {
					logging.LogError(err, "Failed to marshal sync response back to struct")
					time.Sleep(RETRY_CONNECT_DELAY)
				} else if !response.Success {
					logging.LogError(nil, response.Error)
					time.Sleep(RETRY_CONNECT_DELAY)
				} else {
					logging.LogInfo("Successfully synced c2 profile!", "name", c2)
					break
				}
			}
		}
	}
}

package rabbitmq

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/MythicMeta/MythicContainer/custombrowserstructs"
	"github.com/MythicMeta/MythicContainer/logging"
)

type CustomBrowserSyncResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SyncCustomBrowserData(custombrowserDef custombrowserstructs.CustomBrowserDefinition) {
	logging.LogInfo("Syncing consuming container", "name", custombrowserDef.Name)
	response := CustomBrowserSyncResponse{}
	syncMessage := map[string]interface{}{
		"custombrowser":     custombrowserDef,
		"container_version": containerVersion,
	}
	for {
		syncMessageJson, err := json.Marshal(syncMessage)
		if err != nil {
			logging.LogError(err, "Failed to serialize consuming container message to json, %s", err.Error())
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		}
		resp, err := RabbitMQConnection.SendRPCMessage(MYTHIC_EXCHANGE, CUSTOMBROWSER_SYNC_ROUTING_KEY, syncMessageJson, true)
		if err != nil {
			logging.LogError(err, "Failed to send consuming container to Mythic")
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
			logging.LogError(errors.New(response.Error), "waiting and trying again...")
			time.Sleep(RETRY_CONNECT_DELAY)
			continue
		}
		logging.LogInfo("Successfully synced custom browser!", "name", custombrowserDef.Name)
		break
	}
}

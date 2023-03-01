package rabbitmq

import (
	"encoding/json"
	"time"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
)

func SyncPayloadData(syncPayloadName *string) {
	// now make our payloadtype info that we're going to sync
	for _, pt := range agentstructs.AllPayloadData.GetAllPayloadTypeNames() {
		if syncPayloadName == nil || *syncPayloadName == pt {
			logging.LogInfo("Syncing payload type", "name", pt)
			syncMessage := agentstructs.PayloadTypeSyncMessage{}
			response := agentstructs.PayloadTypeSyncMessageResponse{}
			//logging.LogInfo("about to sync over definition", "payload", agentstructs.AllPayloadData.GetPayloadDefinition(), "commands", agentstructs.AllPayloadData.GetCommands())
			syncMessage.PayloadType = agentstructs.AllPayloadData.Get(pt).GetPayloadDefinition()
			syncMessage.CommandList = agentstructs.AllPayloadData.Get(pt).GetCommands()
			agentstructs.AllPayloadData.Get(pt).AddContainerVersion(containerVersion)
			syncMessage.ContainerVersion = agentstructs.AllPayloadData.Get(pt).GetContainerVersion()
			for {
				if syncMessageJson, err := json.Marshal(syncMessage); err != nil {
					logging.LogError(err, "Failed to serialize payload sync message to json, %s", err.Error())
					time.Sleep(RETRY_CONNECT_DELAY)
				} else if resp, err := RabbitMQConnection.SendRPCMessage(MYTHIC_EXCHANGE, PT_SYNC_ROUTING_KEY, syncMessageJson, true); err != nil {
					logging.LogError(err, "Failed to send payload type to Mythic")
					time.Sleep(RETRY_CONNECT_DELAY)
				} else if err := json.Unmarshal(resp, &response); err != nil {
					logging.LogError(err, "Failed to marshal sync response back to struct")
					time.Sleep(RETRY_CONNECT_DELAY)
				} else if !response.Success {
					logging.LogError(nil, response.Error)
					time.Sleep(RETRY_CONNECT_DELAY)
				} else {
					logging.LogInfo("Successfully synced payload type!", "name", pt)
					break
				}
			}
		}
	}
}

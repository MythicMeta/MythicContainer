package rabbitmq

import (
	"encoding/json"
	"errors"
	"time"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
)

func SyncPayloadData(syncPayloadName *string, forcedResync bool) {
	// now make our payloadtype info that we're going to sync
	for _, pt := range agentstructs.AllPayloadData.GetAllPayloadTypeNames() {
		if syncPayloadName == nil || *syncPayloadName == pt {
			logging.LogInfo("Syncing payload type", "name", pt)
			syncMessage := agentstructs.PayloadTypeSyncMessage{
				ForcedResync: forcedResync,
			}
			response := agentstructs.PayloadTypeSyncMessageResponse{}
			//logging.LogInfo("about to sync over definition", "payload", agentstructs.AllPayloadData.GetPayloadDefinition(), "commands", agentstructs.AllPayloadData.GetCommands())
			syncMessage.PayloadType = agentstructs.AllPayloadData.Get(pt).GetPayloadDefinition()
			syncMessage.CommandList = agentstructs.AllPayloadData.Get(pt).GetCommands()
			agentstructs.AllPayloadData.Get(pt).AddContainerVersion(containerVersion)
			syncMessage.ContainerVersion = agentstructs.AllPayloadData.Get(pt).GetContainerVersion()
			for {
				syncMessageJson, err := json.Marshal(syncMessage)
				if err != nil {
					logging.LogError(err, "Failed to serialize payload sync message to json, %s", err.Error())
					time.Sleep(RETRY_CONNECT_DELAY)
					continue
				}
				resp, err := RabbitMQConnection.SendRPCMessage(MYTHIC_EXCHANGE, PT_SYNC_ROUTING_KEY, syncMessageJson, true)
				if err != nil {
					logging.LogError(err, "Failed to send payload type to Mythic")
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
				logging.LogInfo("Successfully synced payload type!", "name", pt)
				break
			}
		}
	}
}

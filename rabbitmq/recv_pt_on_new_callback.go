package rabbitmq

import (
	"encoding/json"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         PT_ON_NEW_CALLBACK,
		RabbitmqProcessingFunction: processPtOnNewCallbackMessages,
	})
}

func processPtOnNewCallbackMessages(msg []byte) {
	incomingMessage := agentstructs.PTOnNewCallbackAllData{}
	response := agentstructs.PTOnNewCallbackResponse{}
	err := json.Unmarshal(msg, &incomingMessage)
	if err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Success = false
		response.Error = "Failed to unmarshal JSON message into structs"
		sendOnNewCallbackResponse(response)
		return
	}
	response.Success = false
	response.AgentCallbackID = incomingMessage.Callback.AgentCallbackID

	onNewCallbackFunc := agentstructs.AllPayloadData.Get(incomingMessage.PayloadType).GetOnNewCallbackFunction()
	if onNewCallbackFunc == nil {
		logging.LogInfo("Failed to get onNewCallbackFunc function. Do you have a function called 'onNewCallbackFunc'? This is an optional function for a payload type to automatically execute tasking and MythicRPC commands when a new callback happens.")
		response.Success = true
	} else {
		response = onNewCallbackFunc(incomingMessage)
	}
	sendOnNewCallbackResponse(response)
	return

}

func sendOnNewCallbackResponse(response agentstructs.PTOnNewCallbackResponse) {
	for {
		err := RabbitMQConnection.SendStructMessage(
			MYTHIC_EXCHANGE,
			PT_ON_NEW_CALLBACK_RESPONSE_ROUTING_KEY,
			"",
			response,
			false,
		)
		if err != nil {
			logging.LogError(err, "Failed to send payload response back to Mythic")
			continue
		}
		return
	}
}

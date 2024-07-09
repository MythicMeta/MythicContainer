package rabbitmq

import (
	"encoding/json"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/mythicutils"
)

func WrapPayloadBuild(msg []byte) {
	//payloadMsg := map[string]interface{}{}
	payloadBuildMsg := agentstructs.PayloadBuildMessage{}
	err := json.Unmarshal(msg, &payloadBuildMsg)
	if err != nil {
		logging.LogError(err, "Failed to process payload build message")
		return
	}
	var payloadBuildResponse agentstructs.PayloadBuildResponse
	payloadBuildFunc := agentstructs.AllPayloadData.Get(payloadBuildMsg.PayloadType).GetBuildFunction()
	if payloadBuildFunc == nil {
		logging.LogError(nil, "Failed to get payload build function. Do you have a function called 'build'?")
		payloadBuildResponse.Success = false
	} else {
		if payloadBuildMsg.WrappedPayloadUUID != nil && *payloadBuildMsg.WrappedPayloadUUID != "" {
			fileContents, err := mythicutils.GetFileFromMythic(*payloadBuildMsg.WrappedPayloadUUID)
			if err != nil {
				payloadBuildResponse.Success = false
				payloadBuildResponse.BuildStdErr = "Failed to get file contents of wrapped payload"
			} else {
				payloadBuildMsg.WrappedPayload = fileContents
				payloadBuildResponse = payloadBuildFunc(payloadBuildMsg)
			}
		} else {
			payloadBuildResponse = payloadBuildFunc(payloadBuildMsg)
		}
	}
	// handle sending off the payload via a web request separately from the rest of the message
	if payloadBuildResponse.Payload != nil {
		if err := mythicutils.SendFileToMythic(payloadBuildResponse.Payload, payloadBuildMsg.PayloadFileUUID); err != nil {
			logging.LogError(err, "Failed to send payload back to Mythic via web request")
			payloadBuildResponse.BuildMessage = payloadBuildResponse.BuildMessage + "\nFailed to send payload back to Mythic: " + err.Error()
			payloadBuildResponse.Success = false
		}
	}
	err = RabbitMQConnection.SendStructMessage(
		MYTHIC_EXCHANGE,
		PT_BUILD_RESPONSE_ROUTING_KEY,
		"",
		payloadBuildResponse,
		false,
	)
	if err != nil {
		logging.LogError(err, "Failed to send payload response back to Mythic")
	}
	logging.LogDebug("Finished processing payload build message")

}

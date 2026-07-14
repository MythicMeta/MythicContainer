package rabbitmq

import (
	"context"
	"encoding/json"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/mythicutils"
)

func WrapPayloadBuild(ctx context.Context, msg []byte) {
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
			directFileToken, err := RequestDirectFileToken(ctx, *payloadBuildMsg.WrappedPayloadUUID, "download")
			if err != nil {
				logging.LogError(err, "Failed to create direct file download token")
				payloadBuildResponse.Success = false
				payloadBuildResponse.BuildStdErr += "\nFailed to get file contents of wrapped payload"
			} else {
				fileContents, err := mythicutils.GetFileFromMythic(ctx, *payloadBuildMsg.WrappedPayloadUUID, directFileToken)
				if err != nil {
					payloadBuildResponse.Success = false
					payloadBuildResponse.BuildStdErr = "Failed to get file contents of wrapped payload"
				} else {
					payloadBuildMsg.WrappedPayload = fileContents
					payloadBuildResponse = payloadBuildFunc(ctx, payloadBuildMsg)
				}
			}
		} else {
			payloadBuildResponse = payloadBuildFunc(ctx, payloadBuildMsg)
		}
	}
	// handle sending off the payload via a web request separately from the rest of the message
	if payloadBuildResponse.Payload != nil {
		directFileToken, err := RequestDirectFileToken(ctx, payloadBuildMsg.PayloadFileUUID, "upload")
		if err == nil {
			err = mythicutils.SendFileToMythic(ctx, payloadBuildResponse.Payload, payloadBuildMsg.PayloadFileUUID, directFileToken)
		}
		if err != nil {
			logging.LogError(err, "Failed to send payload back to Mythic via web request")
			payloadBuildResponse.BuildMessage = payloadBuildResponse.BuildMessage + "\nFailed to send payload back to Mythic: " + err.Error()
			payloadBuildResponse.Success = false
		}
	}
	for {
		err = RabbitMQConnection.SendStructMessageWithContext(ctx,
			MYTHIC_EXCHANGE,
			PT_BUILD_RESPONSE_ROUTING_KEY,
			"",
			payloadBuildResponse,
			false,
		)
		if err != nil {
			logging.LogError(err, "Failed to send payload response back to Mythic")
			continue
		}
		logging.LogDebug("Finished processing payload build message")
		return
	}
}

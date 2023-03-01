package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCPayloadUpdateBuildStepMessage struct {
	PayloadUUID string `json:"payload_uuid"`
	StepName    string `json:"step_name"`
	StepStdout  string `json:"step_stdout"`
	StepStderr  string `json:"step_stderr"`
	StepSuccess bool   `json:"step_success"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCPayloadUpdateBuildStepMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendMythicRPCPayloadUpdateBuildStep(input MythicRPCPayloadUpdateBuildStepMessage) (*MythicRPCPayloadUpdateBuildStepMessageResponse, error) {
	response := MythicRPCPayloadUpdateBuildStepMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_PAYLOAD_UPDATE_BUILD_STEP,
		input,
	); err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	} else if err := json.Unmarshal(responseBytes, &response); err != nil {
		logging.LogError(err, "Failed to parse MythicRPCPayloadUpdateBuildStepMessageResponse response back to struct", "response", response)
		return nil, err
	} else {
		return &response, nil
	}
}

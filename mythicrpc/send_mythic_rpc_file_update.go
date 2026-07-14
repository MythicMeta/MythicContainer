package mythicrpc

import (
	"context"
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
	"github.com/MythicMeta/MythicContainer/utils/mythicutils"
)

type MythicRPCFileUpdateMessage struct {
	AgentFileID      string  `json:"file_id"`
	Comment          string  `json:"comment"`
	Filename         string  `json:"filename"`
	AppendContents   *[]byte `json:"append_contents,omitempty"`
	ReplaceContents  *[]byte `json:"-"`
	Delete           bool    `json:"delete"`
	DeleteAfterFetch *bool   `json:"delete_after_fetch"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCFileUpdateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SendMythicRPCFileUpdate(ctx context.Context, input MythicRPCFileUpdateMessage) (*MythicRPCFileUpdateMessageResponse, error) {
	response := MythicRPCFileUpdateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessageWithContext(
		ctx,
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_FILE_UPDATE,
		input,
	); err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	} else if err := json.Unmarshal(responseBytes, &response); err != nil {
		logging.LogError(err, "Failed to parse response back to struct", "response", response)
		return nil, err
	} else if response.Success {
		if input.ReplaceContents != nil {
			directFileToken, err := getDirectFileToken(ctx, input.AgentFileID, "upload")
			if err != nil {
				response.Success = false
				response.Error = err.Error()
				return &response, nil
			}
			if err := mythicutils.SendFileToMythic(ctx, input.ReplaceContents, input.AgentFileID, directFileToken); err != nil {
				response.Success = false
				response.Error = err.Error()
			}
		}
		return &response, nil
	} else {
		return &response, nil
	}
}

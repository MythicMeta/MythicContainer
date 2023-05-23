package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
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

func SendMythicRPCFileUpdate(input MythicRPCFileUpdateMessage) (*MythicRPCFileUpdateMessageResponse, error) {
	response := MythicRPCFileUpdateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
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
			if err := sendFileToMythic(input.ReplaceContents, input.AgentFileID); err != nil {
				response.Success = false
				response.Error = err.Error()
			}
		}
		return &response, nil
	} else {
		return &response, nil
	}
}

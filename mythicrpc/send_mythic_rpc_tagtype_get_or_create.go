package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCTagTypeGetOrCreateMessage struct {
	TaskID                        int     `json:"task_id"`
	GetOrCreateTagTypeID          *int    `json:"get_or_create_tag_type_id"`
	GetOrCreateTagTypeName        *string `json:"get_or_create_tag_type_name"`
	GetOrCreateTagTypeDescription *string `json:"get_or_create_tag_type_description"`
	GetOrCreateTagTypeColor       *string `json:"get_or_create_tag_type_color"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCTagTypeGetOrCreateMessageResponse struct {
	Success bool                 `json:"success"`
	Error   string               `json:"error"`
	TagType MythicRPCTagTypeData `json:"tagtype"`
}

func SendMythicRPCTagTypeGetOrCreate(input MythicRPCTagTypeGetOrCreateMessage) (*MythicRPCTagTypeGetOrCreateMessageResponse, error) {
	response := MythicRPCTagTypeGetOrCreateMessageResponse{}
	responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_TAGTYPE_GET_OR_CREATE,
		input,
	)
	if err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		logging.LogError(err, "Failed to parse response back to struct", "response", response)
		return nil, err
	}
	return &response, nil
}

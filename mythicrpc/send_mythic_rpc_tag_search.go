package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCTagSearchMessage struct {
	TaskID                int     `json:"task_id"`
	SearchTagID           *int    `json:"search_tag_id"`
	SearchTagTaskID       *int    `json:"search_tag_task_id"`
	SearchTagFileID       *int    `json:"search_tag_file_id,omitempty"`
	SearchTagCredentialID *int    `json:"search_tag_credential_id,omitempty"`
	SearchTagMythicTreeID *int    `json:"search_tag_mythictree_id,omitempty"`
	SearchTagSource       *string `json:"search_tag_source,omitempty"`
	SearchTagData         *string `json:"search_tag_data,omitempty"`
	SearchTagURL          *string `json:"search_tag_url,omitempty"`
}

type MythicRPCTagTypeData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	OperationID int    `json:"operation_id"`
}
type MythicRPCTagData struct {
	ID           int                    `json:"id"`
	TagTypeID    int                    `json:"tagtype_id"`
	TagType      MythicRPCTagTypeData   `json:"tagtype"`
	Data         map[string]interface{} `json:"data"`
	URL          string                 `json:"url"`
	Source       string                 `json:"source"`
	TaskID       *int                   `json:"task_id"`
	FileID       *int                   `json:"file_id"`
	CredentialID *int                   `json:"credential_id"`
	MythicTreeID *int                   `json:"mythic_tree_id"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCTagSearchMessageResponse struct {
	Success bool               `json:"success"`
	Error   string             `json:"error"`
	Tags    []MythicRPCTagData `json:"tags"`
}

func SendMythicRPCTagSearch(input MythicRPCTagSearchMessage) (*MythicRPCTagSearchMessageResponse, error) {
	response := MythicRPCTagSearchMessageResponse{}
	responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_TAG_SEARCH,
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

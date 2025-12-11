package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCCustomBrowserSearchMessage struct {
	TaskID                 *int                             `json:"task_id"`
	OperationID            *int                             `json:"operation_id"`
	GetAllMatchingChildren bool                             `json:"all_matching_children"`
	SearchCustomBrowser    MythicRPCCustomBrowserSearchData `json:"custombrowser"`
}
type MythicRPCCustomBrowserSearchMessageResponse struct {
	Success              bool                                       `json:"success"`
	Error                string                                     `json:"error"`
	CustomBrowserEntries []MythicRPCCustomBrowserSearchDataResponse `json:"custombrowser"`
}
type MythicRPCCustomBrowserSearchData struct {
	TreeType      string      `json:"tree_type" mapstructure:"tree_type"`
	Host          *string     `json:"host" mapstructure:"host"`
	Name          *string     `json:"name" mapstructure:"name"`
	ParentPath    *string     `json:"parent_path" mapstructure:"parent_path"`
	FullPath      *string     `json:"full_path" mapstructure:"full_path"`
	MetadataKey   *string     `json:"metadata_key" mapstructure:"metadata_key"`
	MetadataValue interface{} `json:"metadata_value" mapstructure:"metadata_value"`
	CallbackGroup *string     `json:"callback_group" mapstructure:"callback_group"`
}
type MythicRPCCustomBrowserSearchDataResponse struct {
	TreeType   string                 `json:"tree_type" mapstructure:"tree_type"`
	Host       string                 `json:"host" mapstructure:"host"`
	Name       string                 `json:"name" mapstructure:"name"`
	ParentPath string                 `json:"parent_path" mapstructure:"parent_path"`
	FullPath   string                 `json:"full_path" mapstructure:"full_path"`
	Metadata   map[string]interface{} `json:"metadata" mapstructure:"metadata"`
}

func SendMythicRPCCustomBrowserSearch(input MythicRPCCustomBrowserSearchMessage) (*MythicRPCCustomBrowserSearchMessageResponse, error) {
	response := MythicRPCCustomBrowserSearchMessageResponse{}
	responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.CUSTOMBROWSER_SEARCH,
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

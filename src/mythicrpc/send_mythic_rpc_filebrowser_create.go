package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCFileBrowserCreateMessage struct {
	TaskID      int                                       `json:"task_id"` //required
	FileBrowser MythicRPCFileBrowserCreateFileBrowserData `json:"filebrowser"`
}
type MythicRPCFileBrowserCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type MythicRPCFileBrowserCreateFileBrowserData = agentMessagePostResponseFileBrowser
type agentMessagePostResponseFileBrowser struct {
	Host          string                                         `json:"host" mapstructure:"host"`
	IsFile        bool                                           `json:"is_file" mapstructure:"is_file"`
	Permissions   map[string]interface{}                         `json:"permissions" mapstructure:"permissions"`
	Name          string                                         `json:"name" mapstructure:"name"`
	ParentPath    string                                         `json:"parent_path" mapstructure:"parent_path"`
	Success       bool                                           `json:"success" mapstructure:"success"`
	AccessTime    int64                                          `json:"access_time" mapstructure:"access_time"`
	ModifyTime    int64                                          `json:"modify_time" mapstructure:"modify_time"`
	Size          int64                                          `json:"size" mapstructure:"size"`
	UpdateDeleted *bool                                          `json:"update_deleted,omitempty" mapstructure:"update_deleted,omitempty"` // option to treat this response as full source of truth
	Files         *[]agentMessagePostResponseFileBrowserChildren `json:"files" mapstructure:"files"`
}
type agentMessagePostResponseFileBrowserChildren struct {
	IsFile      bool                   `json:"is_file" mapstructure:"is_file"`
	Permissions map[string]interface{} `json:"permissions" mapstructure:"permissions"`
	Name        string                 `json:"name" mapstructure:"name"`
	AccessTime  int64                  `json:"access_time" mapstructure:"access_time"`
	ModifyTime  int64                  `json:"modify_time" mapstructure:"modify_time"`
	Size        int64                  `json:"size" mapstructure:"size"`
}

func SendMythicRPCFileBrowserCreate(input MythicRPCFileBrowserCreateMessage) (*MythicRPCFileBrowserCreateMessageResponse, error) {
	response := MythicRPCFileBrowserCreateMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_FILEBROWSER_CREATE,
		input,
	); err != nil {
		logging.LogError(err, "Failed to send RPC message")
		return nil, err
	} else if err := json.Unmarshal(responseBytes, &response); err != nil {
		logging.LogError(err, "Failed to parse response back to struct", "response", response)
		return nil, err
	} else {
		return &response, nil
	}
}

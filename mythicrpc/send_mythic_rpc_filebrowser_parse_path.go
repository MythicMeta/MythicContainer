package mythicrpc

import (
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCFileBrowserParsePathMessage struct {
	Path string `json:"path"`
}
type MythicRPCFileBrowserParsePathMessageResponse struct {
	Success      bool         `json:"success"`
	Error        string       `json:"error"`
	AnalyzedPath AnalyzedPath `json:"analyzed_path"`
}
type AnalyzedPath struct {
	PathPieces    []string `json:"path_pieces"`
	PathSeparator string   `json:"path_separator"`
	Host          string   `json:"host"`
}

func SendMythicRPCFileBrowserParsePath(input MythicRPCFileBrowserParsePathMessage) (*MythicRPCFileBrowserParsePathMessageResponse, error) {
	response := MythicRPCFileBrowserParsePathMessageResponse{}
	if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
		rabbitmq.MYTHIC_EXCHANGE,
		rabbitmq.MYTHIC_RPC_FILEBROWSER_PARSE_PATH,
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

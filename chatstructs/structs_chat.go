package chatstructs

import "time"

type ChatContextMessage struct {
	ID                int       `json:"id"`
	AuthorType        string    `json:"author_type"`
	SenderDisplayName string    `json:"sender_display_name"`
	Message           string    `json:"message"`
	CreatedAt         time.Time `json:"created_at"`
}

type ChatRequestMessage struct {
	ContainerName     string                 `json:"container_name"`
	OperationID       int                    `json:"operation_id"`
	ChannelID         int                    `json:"channel_id"`
	APITokenID        int                    `json:"apitokens_id"`
	ChannelName       string                 `json:"channel_name"`
	ChannelSlug       string                 `json:"channel_slug"`
	RequestID         int                    `json:"request_id"`
	RequestMessageID  int                    `json:"request_message_id"`
	ResponseMessageID int                    `json:"response_message_id"`
	Model             string                 `json:"model"`
	Prompt            string                 `json:"prompt"`
	Config            map[string]interface{} `json:"config"`
	Context           []ChatContextMessage   `json:"context"`
	Secrets           map[string]interface{} `json:"secrets"`
}

type ChatResponseMessage struct {
	OperationID       int                    `json:"operation_id"`
	RequestID         int                    `json:"request_id"`
	ResponseMessageID int                    `json:"response_message_id"`
	Content           string                 `json:"content"`
	IsDelta           bool                   `json:"is_delta"`
	Complete          bool                   `json:"complete"`
	Status            string                 `json:"status"`
	Error             string                 `json:"error"`
	Metadata          map[string]interface{} `json:"metadata"`
}

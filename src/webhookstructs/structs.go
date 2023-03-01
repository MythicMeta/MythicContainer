package webhookstructs

type webhookMessageBase struct {
	OperationID      int          `json:"operation_id"`
	OperationName    string       `json:"operation_name"`
	OperationWebhook string       `json:"operation_webhook"`
	OperationChannel string       `json:"operation_channel"`
	OperatorUsername string       `json:"operator_username"`
	Action           WEBHOOK_TYPE `json:"action"`
}

type SlackWebhookMessage struct {
	Channel     string                          `json:"channel"`
	Username    string                          `json:"username"`
	IconEmoji   string                          `json:"icon_emoji"`
	Attachments []SlackWebhookMessageAttachment `json:"attachments"`
}

type SlackWebhookMessageAttachment struct {
	Title  string                                `json:"fallback"`
	Color  string                                `json:"color,omitempty"`
	Blocks *[]SlackWebhookMessageAttachmentBlock `json:"blocks,omitempty"`
}

type SlackWebhookMessageAttachmentBlock struct {
	Type   string                                    `json:"type"`
	Text   *SlackWebhookMessageAttachmentBlockText   `json:"text,omitempty"`
	Fields *[]SlackWebhookMessageAttachmentBlockText `json:"fields,omitempty"`
}

type SlackWebhookMessageAttachmentBlockText struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

func GetNewDefaultWebhookMessage() SlackWebhookMessage {
	headerBlockText := SlackWebhookMessageAttachmentBlockText{
		Type: "mrkdwn",
		Text: "You have a new message from Mythic!",
	}

	blocks := []SlackWebhookMessageAttachmentBlock{
		{
			Type: "section",
			Text: &headerBlockText,
		},
		{
			Type: "divider",
		},
	}
	newMsg := SlackWebhookMessage{
		Username:  "Mythic",
		IconEmoji: ":mythic:",
		Channel:   "#mythic-notifications",
		Attachments: []SlackWebhookMessageAttachment{
			{
				Title:  "New Mythic Message!",
				Color:  "#b366ff",
				Blocks: &blocks,
			},
		},
	}
	return newMsg
}

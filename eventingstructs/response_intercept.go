package eventingstructs

type ResponseInterceptMessage struct {
	EventStepInstanceID int                    `json:"eventstepinstance_id"`
	ResponseID          int                    `json:"response_id"`
	CallbackID          int                    `json:"callback_id"`
	CallbackDisplayID   int                    `json:"callback_display_id"`
	AgentCallbackID     string                 `json:"agent_callback_id"`
	ContainerName       string                 `json:"container_name"`
	Environment         map[string]interface{} `json:"environment"`
	Inputs              map[string]interface{} `json:"inputs"`
	ActionData          map[string]interface{} `json:"action_data"`
}

type ResponseInterceptMessageResponse struct {
	EventStepInstanceID int                    `json:"eventstepinstance_id" mapstructure:"eventstepinstance_id"`
	ResponseID          int                    `json:"response_id" mapstructure:"response_id"`
	Success             bool                   `json:"success" mapstructure:"success"`
	StdOut              string                 `json:"stdout" mapstructure:"stdout"`
	StdErr              string                 `json:"stderr" mapstructure:"stderr"`
	Outputs             map[string]interface{} `json:"outputs" mapstructure:"outputs"`
	Response            string                 `json:"response" mapstructure:"response"`
}

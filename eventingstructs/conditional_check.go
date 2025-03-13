package eventingstructs

type ConditionalCheckEventingMessage struct {
	EventStepInstanceID int                    `json:"eventstepinstance_id"`
	FunctionName        string                 `json:"function_name"`
	ContainerName       string                 `json:"container_name"`
	Environment         map[string]interface{} `json:"environment"`
	Inputs              map[string]interface{} `json:"inputs"`
	ActionData          map[string]interface{} `json:"action_data"`
}
type ConditionalCheckEventingMessageResponse struct {
	EventStepInstanceID int                    `json:"eventstepinstance_id" mapstructure:"eventstepinstance_id"`
	Success             bool                   `json:"success" mapstructure:"success"`
	StdOut              string                 `json:"stdout" mapstructure:"stdout"`
	StdErr              string                 `json:"stderr" mapstructure:"stderr"`
	Outputs             map[string]interface{} `json:"outputs" mapstructure:"outputs"`
	SkipStep            bool                   `json:"skip_step" mapstructure:"skip_step"`
}

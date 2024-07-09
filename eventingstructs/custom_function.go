package eventingstructs

type NewCustomEventingMessage struct {
	EventStepInstanceID int                    `json:"eventstepinstance_id"`
	FunctionName        string                 `json:"function_name"`
	ContainerName       string                 `json:"container_name"`
	Environment         map[string]interface{} `json:"environment"`
	Inputs              map[string]interface{} `json:"inputs"`
	ActionData          map[string]interface{} `json:"action_data"`
}
type NewCustomEventingMessageResponse struct {
	EventStepInstanceID int                    `json:"eventstepinstance_id"`
	Success             bool                   `json:"success"`
	StdOut              string                 `json:"stdout"`
	StdErr              string                 `json:"stderr"`
	Outputs             map[string]interface{} `json:"outputs"`
}

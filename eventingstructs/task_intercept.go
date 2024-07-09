package eventingstructs

type TaskInterceptMessage struct {
	EventStepInstanceID int                    `json:"eventstepinstance_id"`
	TaskID              int                    `json:"task_id"`
	CallbackID          int                    `json:"callback_id"`
	ContainerName       string                 `json:"container_name"`
	Environment         map[string]interface{} `json:"environment"`
	Inputs              map[string]interface{} `json:"inputs"`
	ActionData          map[string]interface{} `json:"action_data"`
}

type TaskInterceptMessageResponse struct {
	EventStepInstanceID int                    `json:"eventstepinstance_id" mapstructure:"eventstepinstance_id"`
	TaskID              int                    `json:"task_id" mapstructure:"task_id"`
	Success             bool                   `json:"success" mapstructure:"success"`
	StdOut              string                 `json:"stdout" mapstructure:"stdout"`
	StdErr              string                 `json:"stderr" mapstructure:"stderr"`
	BlockTask           bool                   `json:"block_task" mapstructure:"block_task"`
	BypassRole          OPSEC_ROLE             `json:"bypass_role" mapstructure:"bypass_role"`
	BypassMessage       string                 `json:"bypass_message" mapstructure:"bypass_message"`
	Outputs             map[string]interface{} `json:"outputs" mapstructure:"outputs"`
}
type OPSEC_ROLE string

const (
	OPSEC_ROLE_LEAD           OPSEC_ROLE = "lead"
	OPSEC_ROLE_OPERATOR                  = "operator"
	OPSEC_ROLE_OTHER_OPERATOR            = "other_operator"
)

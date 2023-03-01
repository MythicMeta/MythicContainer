package mythicrpc

type PTTaskMessageTaskData struct {
	ID                                 int    `json:"id"`
	DisplayID                          int    `json:"display_id"`
	AgentTaskID                        string `json:"agent_task_id"`
	CommandName                        string `json:"command_name"`
	Params                             string `json:"params"`
	Timestamp                          string `json:"timestamp"`
	CallbackID                         int    `json:"callback_id"`
	Status                             string `json:"status"`
	OriginalParams                     string `json:"original_params"`
	DisplayParams                      string `json:"display_params"`
	Comment                            string `json:"comment"`
	Stdout                             string `json:"stdout"`
	Stderr                             string `json:"stderr"`
	Completed                          bool   `json:"completed"`
	OpsecPreBlocked                    bool   `json:"opsec_pre_blocked"`
	OpsecPreMessage                    string `json:"opsec_pre_message"`
	OpsecPreBypassed                   bool   `json:"opsec_pre_bypassed"`
	OpsecPreBypassRole                 string `json:"opsec_pre_bypass_role"`
	OpsecPostBlocked                   bool   `json:"opsec_post_blocked"`
	OpsecPostMessage                   string `json:"opsec_post_message"`
	OpsecPostBypassed                  bool   `json:"opsec_post_bypassed"`
	OpsecPostBypassRole                string `json:"opsec_post_bypass_role"`
	ParentTaskID                       int    `json:"parent_task_id"`
	OperatorUsername                   string `json:"operator_username"`
	SubtaskCallbackFunction            string `json:"subtask_callback_function"`
	SubtaskCallbackFunctionCompleted   bool   `json:"subtask_callback_function_completed"`
	GroupCallbackFunction              string `json:"group_callback_function"`
	GroupCallbackFunctionCompleted     bool   `json:"group_callback_function_completed"`
	CompletedCallbackFunction          string `json:"completed_callback_function"`
	CompletedCallbackFunctionCompleted bool   `json:"completed_callback_function_completed"`
	SubtaskGroupName                   string `json:"subtask_group_name"`
	TaskingLocation                    string `json:"tasking_location"`
	ParameterGroupName                 string `json:"parameter_group_name"`
	TokenID                            int    `json:"token_id"`
}

// exporting a payload configuration
type PayloadConfiguration struct {
	Description        string                                `json:"description"`
	PayloadType        string                                `json:"payload_type" binding:"required"`
	C2Profiles         *[]PayloadConfigurationC2Profile      `json:"c2_profiles,omitempty"`
	BuildParameters    *[]PayloadConfigurationBuildParameter `json:"build_parameters,omitempty"`
	Commands           []string                              `json:"commands"`
	SelectedOS         string                                `json:"selected_os" binding:"required"`
	Filename           string                                `json:"filename" binding:"required"`
	WrappedPayloadUUID string                                `json:"wrapped_payload"`
	UUID               string                                `json:"uuid"`
	AgentFileID        string                                `json:"agent_file_id"`
}
type PayloadConfigurationC2Profile struct {
	Name       string                 `json:"c2_profile"`
	Parameters map[string]interface{} `json:"c2_profile_parameters"`
}
type PayloadConfigurationBuildParameter struct {
	Name  string      `json:"name" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

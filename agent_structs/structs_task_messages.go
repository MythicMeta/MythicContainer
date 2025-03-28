package agentstructs

// PT_TASK_* structs

type PTTaskMessageAllData struct {
	// Task - Read-only data about the task
	Task PTTaskMessageTaskData `json:"task"`
	// Callback - Read-only data about the callback
	Callback PTTaskMessageCallbackData `json:"callback"`
	// BuildParameters - Read-only data about the build parameters
	BuildParameters []PayloadConfigurationBuildParameter `json:"build_parameters"`
	// Commands - Read-only data about the commands built into the callback
	Commands []string `json:"commands"`
	// Payload - Read-only data about the backing payload for this task
	Payload PTTaskMessagePayloadData `json:"payload"`
	// C2Profiles - Read-only data about the c2 profiles and their values for this callback
	C2Profiles []PayloadConfigurationC2Profile `json:"c2info"`
	// PayloadType - Read-only the name of the payload type associated with this callback
	PayloadType string `json:"payload_type"`
	// CommandPayloadType The name of the payload type associated with this task
	CommandPayloadType string `json:"command_payload_type"`
	// Secrets - Map of user supplied secrets to their values to help with tasking
	Secrets map[string]interface{} `json:"secrets"`
	// Args - Read-Write argument data for adding/removing/modifying args associated with this task instance.
	// Mainly for create tasking function to augment parameters
	Args PTTaskMessageArgsData
}

type PTTaskMessageTaskData struct {
	ID                                 int    `json:"id"`
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
	OperatorUsername                   string `json:"operator_username"`
	OperatorID                         int    `json:"operator_id"`
	OpsecPreBlocked                    bool   `json:"opsec_pre_blocked"`
	OpsecPreMessage                    string `json:"opsec_pre_message"`
	OpsecPreBypassed                   bool   `json:"opsec_pre_bypassed"`
	OpsecPreBypassRole                 string `json:"opsec_pre_bypass_role"`
	OpsecPostBlocked                   bool   `json:"opsec_post_blocked"`
	OpsecPostMessage                   string `json:"opsec_post_message"`
	OpsecPostBypassed                  bool   `json:"opsec_post_bypassed"`
	OpsecPostBypassRole                string `json:"opsec_post_bypass_role"`
	ParentTaskID                       int    `json:"parent_task_id"`
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
	IsInteractiveTask                  bool   `json:"is_interactive_task"`
	InteractiveTaskType                int    `json:"interactive_task_type"`
}

type PTTaskMessageCallbackData struct {
	ID                  int      `json:"id"`
	DisplayID           int      `json:"display_id"`
	AgentCallbackID     string   `json:"agent_callback_id"`
	InitCallback        string   `json:"init_callback"`
	LastCheckin         string   `json:"last_checkin"`
	User                string   `json:"user"`
	Host                string   `json:"host"`
	PID                 int      `json:"pid"`
	IP                  string   `json:"ip"`
	IPs                 []string `json:"ips"`
	ExternalIp          string   `json:"external_ip"`
	ProcessName         string   `json:"process_name"`
	Description         string   `json:"description"`
	OperatorID          int      `json:"operator_id"`
	OperatorUsername    string   `json:"operator_username"`
	Active              bool     `json:"active"`
	RegisteredPayloadID int      `json:"registered_payload_id"`
	IntegrityLevel      int      `json:"integrity_level"`
	Locked              bool     `json:"locked"`
	OperationID         int      `json:"operation_id"`
	OperationName       string   `json:"operation_name"`
	CryptoType          string   `json:"crypto_type"`
	DecKey              []byte   `json:"dec_key"`
	EncKey              []byte   `json:"enc_key"`
	OS                  string   `json:"os"`
	Architecture        string   `json:"architecture"`
	Domain              string   `json:"domain"`
	ExtraInfo           string   `json:"extra_info"`
	SleepInfo           string   `json:"sleep_info"`
}

type PTTaskMessagePayloadData struct {
	OS          string `json:"os"`
	UUID        string `json:"uuid"`
	PayloadType string `json:"payload_type"`
}

type PT_TASK_FUNCTION_STATUS = string

type PtTaskFunctionParseArgString func(args *PTTaskMessageArgsData, input string) error
type PtTaskFunctionParseArgDictionary func(args *PTTaskMessageArgsData, input map[string]interface{}) error

// PTTaskMessageArgsData - struct for tracking, adding, removing, updating, validating, etc arguments for a task.
// If you want to set your own manual arguments, use the .SetManualArgs function.
type PTTaskMessageArgsData struct {
	args            []CommandParameter
	commandLine     string
	rawCommandLine  string
	taskingLocation string
	manualArgs      *string
	// manualParameterGroupName use this in case of a user-explicit group
	manualParameterGroupName string
	// initialParameterGroupName use this in case of multiple matching parameter groups
	initialParameterGroupName string
}

const (
	PT_TASK_FUNCTION_STATUS_OPSEC_PRE                        PT_TASK_FUNCTION_STATUS = "OPSEC Pre Check Running..."
	PT_TASK_FUNCTION_STATUS_OPSEC_PRE_ERROR                                          = "Error: opsec check - check task stdout/stderr"
	PT_TASK_FUNCTION_STATUS_OPSEC_PRE_BLOCKED                                        = "OPSEC Pre Blocked"
	PT_TASK_FUNCTION_STATUS_PREPROCESSING                                            = "creating task..."
	PT_TASK_FUNCTION_STATUS_PREPROCESSING_ERROR                                      = "Error: creating task - check task stdout/stderr"
	PT_TASK_FUNCTION_STATUS_OPSEC_POST                                               = "OPSEC Post Check Running..."
	PT_TASK_FUNCTION_STATUS_OPSEC_POST_ERROR                                         = "Error: opsec check - check task stdout/stderr"
	PT_TASK_FUNCTION_STATUS_OPSEC_POST_BLOCKED                                       = "OPSEC Post Blocked"
	PT_TASK_FUNCTION_STATUS_SUBMITTED                                                = "submitted"
	PT_TASK_FUNCTION_STATUS_PROCESSING                                               = "agent processing"
	PT_TASK_FUNCTION_STATUS_DELEGATING                                               = "delegating tasks..."
	PT_TASK_FUNCTION_STATUS_COMPLETION_FUNCTION                                      = "Completion Function Running..."
	PT_TASK_FUNCTION_STATUS_COMPLETION_FUNCTION_ERROR                                = "Error: completion function - check task stdout/stderr"
	PT_TASK_FUNCTION_STATUS_SUBTASK_COMPLETED_FUNCTION                               = "SubTask Completion Function Running..."
	PT_TASK_FUNCTION_STATUS_SUBTASK_COMPLETED_FUNCTION_ERROR                         = "Error: subtask completion function - check task stdout/stderr"
	PT_TASK_FUNCTION_STATUS_GROUP_COMPLETED_FUNCTION                                 = "Group Completion Function Running..."
	PT_TASK_FUNCTION_STATUS_GROUP_COMPLETED_FUNCTION_ERROR                           = "Error: group completion function - check task stdout/stderr"
	PT_TASK_FUNCTION_STATUS_COMPLETED                                                = "completed"
	PT_TASK_FUNCTION_STATUS_PROCESSED                                                = "processed, waiting for more messages..."
)

// Tasking step 1:
// Task message/process before running create_tasking function
//
//	opportunity to run any necessary opsec checks/blocks before the logic in create_tasking runs
//		which can spawn subtasks outside of the opsec checks
type OPSEC_ROLE string

const (
	OPSEC_ROLE_LEAD           OPSEC_ROLE = "lead"
	OPSEC_ROLE_OPERATOR                  = "operator"
	OPSEC_ROLE_OTHER_OPERATOR            = "other_operator"
)

type PtTaskFunctionOPSECPre func(*PTTaskMessageAllData) PTTTaskOPSECPreTaskMessageResponse
type PTTTaskOPSECPreTaskMessageResponse struct {
	TaskID             int        `json:"task_id"`
	Success            bool       `json:"success"`
	Error              string     `json:"error"`
	OpsecPreBlocked    bool       `json:"opsec_pre_blocked"`
	OpsecPreMessage    string     `json:"opsec_pre_message"`
	OpsecPreBypassed   *bool      `json:"opsec_pre_bypassed,omitempty"`
	OpsecPreBypassRole OPSEC_ROLE `json:"opsec_pre_bypass_role"`
}

// Tasking step 2:
// Task message/process to run the create_tasking function

// PtTaskFunctionCreateTasking - Process the tasking request from the user. If you want to access/modify the arguments
// for this task, use the Task.Args.* functions.
type PtTaskFunctionCreateTasking func(*PTTaskMessageAllData) PTTaskCreateTaskingMessageResponse
type PTTaskCreateTaskingMessageResponse struct {
	// TaskID - the task associated with the create tasking function - this will be automatically filled in for you
	TaskID int `json:"task_id"`
	// Success - indicate if the create tasking function ran successfully or not
	Success bool `json:"success"`
	// Error - if you want to provide an error message about some error you hit while executing the create tasking
	Error string `json:"error"`
	// CommandName - if you want to change the associated command name that's sent down to the agent
	CommandName *string `json:"command_name,omitempty"`
	// TaskStatus - if you want to manually set the task status to be something other than default
	TaskStatus *string `json:"task_status,omitempty"`
	// DisplayParams - if you want to change the display parameters for your task to be something other than the default JSON
	DisplayParams *string `json:"display_params,omitempty"`
	// Stdout - Provide any task-based stdout
	Stdout *string `json:"stdout,omitempty"`
	// Stderr - Provide any task-based stderr
	Stderr *string `json:"stderr,omitempty"`
	// Completed - identify if the task is already completed and shouldn't be sent down to the agent
	Completed *bool `json:"completed,omitempty"`
	// TokenID - identifier for the token id associated with this task - normally doesn't need to be set unless you're changing it
	TokenID *uint64 `json:"token_id,omitempty"`
	// CompletionFunctionName - name of the completion function to call from the Command's TaskCompletionFunctions dictionary
	CompletionFunctionName *string `json:"completion_function_name,omitempty"`
	// ParameterGroupName - Don't set this explicitly. If you want to set the name of the parameter group explicitly, use
	// the taskData.Args.SetManualParameterGroup("name here") function.
	ParameterGroupName string `json:"parameter_group_name"`
	// ReprocessAtNewCommandPayloadType - the name of the current payload type or payload type associated with an updated CommandName field for execution to then go to for further processing
	ReprocessAtNewCommandPayloadType string `json:"reprocess_at_new_command_payload_type"`
}

// Tasking step 3:
// Task message/process after running create_tasking but before the task can be picked up by an agent
//
//	this is the time to check any artifacts generated from create_tasking
type PtTaskFunctionOPSECPost func(*PTTaskMessageAllData) PTTaskOPSECPostTaskMessageResponse
type PTTaskOPSECPostTaskMessageResponse struct {
	TaskID              int        `json:"task_id"`
	Success             bool       `json:"success"`
	Error               string     `json:"error"`
	OpsecPostBlocked    bool       `json:"opsec_post_blocked"`
	OpsecPostMessage    string     `json:"opsec_post_message"`
	OpsecPostBypassed   *bool      `json:"opsec_post_bypassed,omitempty"`
	OpsecPostBypassRole OPSEC_ROLE `json:"opsec_post_bypass_role"`
}

// Tasking step 4:
// Run this when the specified task completes
type SubtaskGroupName = string

type PTTaskCompletionFunctionMessage struct {
	TaskData               *PTTaskMessageAllData `json:"task"`
	SubtaskData            *PTTaskMessageAllData `json:"subtask,omitempty"`
	SubtaskGroup           *SubtaskGroupName     `json:"subtask_group_name,omitempty"`
	CompletionFunctionName string                `json:"function_name"`
}

// PTTaskCompletionFunction takes in taskData, subtaskData, groupName
// taskData is always your current task
// subtaskData is optional if this is executing once a subtask finishes execution
// subtaskGroupName is optional if the subtask was part of a named group
type PTTaskCompletionFunction func(*PTTaskMessageAllData, *PTTaskMessageAllData, *SubtaskGroupName) PTTaskCompletionFunctionMessageResponse
type PTTaskCompletionFunctionMessageResponse struct {
	TaskID                 int     `json:"task_id"`
	ParentTaskId           int     `json:"parent_task_id"`
	Success                bool    `json:"success"`
	Error                  string  `json:"error"`
	TaskStatus             *string `json:"task_status,omitempty"`
	DisplayParams          *string `json:"display_params,omitempty"`
	Stdout                 *string `json:"stdout,omitempty"`
	Stderr                 *string `json:"stderr,omitempty"`
	Completed              *bool   `json:"completed,omitempty"`
	TokenID                *int    `json:"token_id,omitempty"`
	CompletionFunctionName *string `json:"completion_function_name,omitempty"`
	Params                 *string `json:"params,omitempty"`
	ParameterGroupName     *string `json:"parameter_group_name,omitempty"`
}

// Tasking step 5:
// Task message/process to run for more manual processing of a message's response data
type PtTaskProcessResponseMessage struct {
	TaskData *PTTaskMessageAllData `json:"task"`
	Response interface{}           `json:"response"`
}
type PtTaskFunctionProcessResponse func(PtTaskProcessResponseMessage) PTTaskProcessResponseMessageResponse
type PTTaskProcessResponseMessageResponse struct {
	TaskID  int    `json:"task_id"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// On New Callback Structs

type PTOnNewCallbackAllData struct {
	Callback        PTTaskMessageCallbackData            `json:"callback"`
	BuildParameters []PayloadConfigurationBuildParameter `json:"build_parameters"`
	Commands        []string                             `json:"commands"`
	Payload         PTTaskMessagePayloadData             `json:"payload"`
	C2Profiles      []PayloadConfigurationC2Profile      `json:"c2info"`
	PayloadType     string                               `json:"payload_type"`
	Secrets         map[string]interface{}               `json:"secrets"`
}

type PTOnNewCallbackResponse struct {
	AgentCallbackID string `json:"agent_callback_id"`
	Success         bool   `json:"success"`
	Error           string `json:"error"`
}

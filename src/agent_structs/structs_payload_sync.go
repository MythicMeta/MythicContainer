package agentstructs

import (
	"encoding/json"
)

// PAYLOAD_SYNC STRUCTS
type PayloadTypeSyncMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type PayloadTypeSyncMessage struct {
	PayloadType      PayloadType `json:"payload_type"`
	CommandList      []Command   `json:"commands"`
	ContainerVersion string      `json:"container_version"`
}

type BuildParameterType = string

const (
	BUILD_PARAMETER_TYPE_STRING          BuildParameterType = "String"
	BUILD_PARAMETER_TYPE_BOOLEAN                            = "Boolean"
	BUILD_PARAMETER_TYPE_CHOOSE_ONE                         = "ChooseOne"
	BUILD_PARAMETER_TYPE_CHOOSE_MULTIPLE                    = "ChooseMultiple"
	BUILD_PARAMETER_TYPE_DATE                               = "Date"
	BUILD_PARAMETER_TYPE_DICTIONARY                         = "Dictionary"
	BUILD_PARAMETER_TYPE_ARRAY                              = "Array"
	BUILD_PARAMETER_TYPE_NUMBER                             = "Number"
)

type BuildParameter struct {
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Required          bool                       `json:"required"`
	VerifierRegex     string                     `json:"verifier_regex"`
	DefaultValue      interface{}                `json:"default_value"`
	ParameterType     BuildParameterType         `json:"parameter_type"`
	FormatString      string                     `json:"format_string"`
	Randomize         bool                       `json:"randomize"`
	IsCryptoType      bool                       `json:"crypto_type"`
	Choices           []string                   `json:"choices"`
	DictionaryChoices []BuildParameterDictionary `json:"dictionary_choices"`
}

type BuildStep struct {
	StepName        string `json:"step_name"`
	StepDescription string `json:"step_description"`
}

type PTRPCOtherServiceRPCMessage struct {
	ServiceName                 string                 `json:"service_name"` //required
	ServiceRPCFunction          string                 `json:"service_function"`
	ServiceRPCFunctionArguments map[string]interface{} `json:"service_arguments"`
}
type PTRPCOtherServiceRPCMessageResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Result  map[string]interface{} `json:"result"`
}

type PayloadType struct {
	Name                                   string                                                                                   `json:"name"`
	FileExtension                          string                                                                                   `json:"file_extension"`
	Author                                 string                                                                                   `json:"author"`
	SupportedOS                            []string                                                                                 `json:"supported_os"`
	Wrapper                                bool                                                                                     `json:"wrapper"`
	CanBeWrappedByTheFollowingPayloadTypes []string                                                                                 `json:"supported_wrapper_payload_types"`
	SupportsDynamicLoading                 bool                                                                                     `json:"supports_dynamic_load"`
	Description                            string                                                                                   `json:"description"`
	SupportedC2Profiles                    []string                                                                                 `json:"supported_c2_profiles"`
	TranslationContainerName               string                                                                                   `json:"translation_container_name"`
	MythicEncryptsData                     bool                                                                                     `json:"mythic_encrypts"`
	BuildParameters                        []BuildParameter                                                                         `json:"build_parameters"`
	BuildSteps                             []BuildStep                                                                              `json:"build_steps"`
	AgentIcon                              *[]byte                                                                                  `json:"agent_icon"` // automatically filled in based on Name
	CustomRPCFunctions                     map[string]func(message PTRPCOtherServiceRPCMessage) PTRPCOtherServiceRPCMessageResponse `json:"-"`
}

type Command struct {
	Name                           string                              `json:"name"`
	NeedsAdminPermissions          bool                                `json:"needs_admin_permissions"`
	HelpString                     string                              `json:"help_string"`
	Description                    string                              `json:"description"`
	Version                        uint32                              `json:"version"`
	SupportedUIFeatures            []string                            `json:"supported_ui_features"`
	Author                         string                              `json:"author"`
	MitreAttackMappings            []string                            `json:"attack"`
	ScriptOnlyCommand              bool                                `json:"script_only"`
	CommandAttributes              CommandAttribute                    `json:"attributes"`
	CommandParameters              []CommandParameter                  `json:"parameters"`
	AssociatedBrowserScript        *BrowserScript                      `json:"browserscript,omitempty"`
	TaskFunctionOPSECPre           PtTaskFunctionOPSECPre              `json:"-"`
	TaskFunctionCreateTasking      PtTaskFunctionCreateTasking         `json:"-"`
	TaskFunctionProcessResponse    PtTaskFunctionProcessResponse       `json:"-"`
	TaskFunctionOPSECPost          PtTaskFunctionOPSECPost             `json:"-"`
	TaskFunctionParseArgString     PtTaskFunctionParseArgString        `json:"-"`
	TaskFunctionParseArgDictionary PtTaskFunctionParseArgDictionary    `json:"-"`
	TaskCompletionFunctions        map[string]PTTaskCompletionFunction `json:"-"`
}
type CommandParameterType = string

const (
	COMMAND_PARAMETER_TYPE_STRING          CommandParameterType = "String"
	COMMAND_PARAMETER_TYPE_BOOLEAN                              = "Boolean"
	COMMAND_PARAMETER_TYPE_CHOOSE_ONE                           = "ChooseOne"
	COMMAND_PARAMETER_TYPE_CHOOSE_MULTIPLE                      = "ChooseMultiple"
	COMMAND_PARAMETER_TYPE_FILE                                 = "File"
	COMMAND_PARAMETER_TYPE_ARRAY                                = "Array"
	COMMAND_PARAMETER_TYPE_CREDENTIAL                           = "CredentialJson"
	COMMAND_PARAMETER_TYPE_NUMBER                               = "Number"
	COMMAND_PARAMETER_TYPE_PAYLOAD_LIST                         = "PayloadList"
	COMMAND_PARAMETER_TYPE_CONNECTION_INFO                      = "AgentConnect"
	COMMAND_PARAMETER_TYPE_LINK_INFO                            = "LinkInfo"
)

type CommandParameter struct {
	Name                                    string                        `json:"name"`
	ModalDisplayName                        string                        `json:"display_name"`
	CLIName                                 string                        `json:"cli_name"`
	ParameterType                           CommandParameterType          `json:"parameter_type"`
	Description                             string                        `json:"description"`
	Choices                                 []string                      `json:"choices"`
	DefaultValue                            interface{}                   `json:"default_value"`
	SupportedAgents                         []string                      `json:"supported_agents"`
	SupportedAgentBuildParameters           map[string]string             `json:"supported_agent_build_parameters"`
	ChoicesAreAllCommands                   bool                          `json:"choices_are_all_commands"`
	ChoicesAreLoadedCommands                bool                          `json:"choices_are_loaded_commands"`
	FilterCommandChoicesByCommandAttributes map[string]string             `json:"choice_filter_by_command_attributes"`
	DynamicQueryFunction                    PTTaskingDynamicQueryFunction `json:"dynamic_query_function"`
	ParameterGroupInformation               []ParameterGroupInfo          `json:"parameter_group_info"`
	value                                   interface{}                   // the current value for the parameter
	userSupplied                            bool                          // was this value supplied by the user or a default value
}

type PTTaskingDynamicQueryFunction func(PTRPCDynamicQueryFunctionMessage) []string

func (f PTTaskingDynamicQueryFunction) MarshalJSON() ([]byte, error) {
	if f != nil {
		return json.Marshal("foo")
	} else {
		return json.Marshal("")
	}
}

type CommandAttribute struct {
	CommandIsInjectableIntoProcess                  bool              `json:"spawn_and_injectable"`
	SupportedOS                                     []string          `json:"supported_os"`
	CommandIsBuiltin                                bool              `json:"builtin"`
	CommandIsSuggested                              bool              `json:"suggested_command"`
	CommandCanOnlyBeLoadedLater                     bool              `json:"load_only"`
	FilterCommandAvailabilityByAgentBuildParameters map[string]string `json:"filter_by_build_parameter"`
	AdditionalAttributes                            map[string]string `json:"additional_items"`
}

type ParameterGroupInfo struct {
	ParameterIsRequired   bool              `json:"required"`
	GroupName             string            `json:"group_name"`
	UIModalPosition       uint32            `json:"ui_position"`
	AdditionalInformation map[string]string `json:"additional_info"`
}

type BrowserScript struct {
	ScriptPath     string `json:"-"`
	Author         string `json:"author"`
	ScriptContents string `json:"script"`
}
type C2ParameterDictionary struct {
	Name         string `json:"name"`
	DefaultValue string `json:"default_value"`
	DefaultShow  bool   `json:"default_show"`
}

type BuildParameterDictionary C2ParameterDictionary

// building just an ad-hoc c2 profile for an already existing payload
type PayloadBuildC2ProfileMessage struct {
	PayloadUUID     string                 `json:"uuid"`
	Parameters      map[string]interface{} `json:"parameters"`
	BuildParameters map[string]interface{} `json:"build_parameters"`
	SelectedOS      string                 `json:"selected_os"`
	PayloadType     string                 `json:"payload_type"`
}

type PayloadBuildC2ProfileMessageResponse struct {
	PayloadUUID  string  `json:"uuid"`
	Status       string  `json:"status"`
	Payload      *[]byte `json:"payload,omitempty"`
	BuildStdErr  string  `json:"build_stderr"`
	BuildStdOut  string  `json:"build_stdout"`
	BuildMessage string  `json:"build_message"`
}

// exporting a payload configuration
type PayloadConfiguration struct {
	Description        string                                `json:"tag"`
	PayloadType        string                                `json:"payload_type" binding:"required"`
	C2Profiles         *[]PayloadConfigurationC2Profile      `json:"c2_profiles,omitempty"`
	BuildParameters    *[]PayloadConfigurationBuildParameter `json:"build_parameters,omitempty"`
	Commands           []string                              `json:"commands"`
	SelectedOS         string                                `json:"selected_os" binding:"required"`
	Filename           string                                `json:"filename" binding:"required"`
	WrappedPayloadUUID string                                `json:"wrapped_payload"`
}
type PayloadConfigurationC2Profile struct {
	Name       string                 `json:"c2_profile"`
	Parameters map[string]interface{} `json:"c2_profile_parameters"`
}
type PayloadConfigurationBuildParameter struct {
	Name  string      `json:"name" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

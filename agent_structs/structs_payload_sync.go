package agentstructs

import (
	"encoding/json"
)

// PayloadTypeSyncMessageResponse - A message back from Mythic indicating if the Payload Sync was successful or not
type PayloadTypeSyncMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// PayloadTypeSyncMessage - A sync message to Mythic describing this Payload Type
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
	BUILD_PARAMETER_TYPE_FILE                               = "File"
	BUILD_PARAMETER_TYPE_TYPED_ARRAY                        = "TypedArray"
)

// BuildParameter - A structure defining the metadata about a build parameter for the user to select when building a payload.
type BuildParameter struct {
	// Name - the name of the build parameter for use during the Payload Type's build function
	Name string `json:"name"`
	// Description - the description of the build parameter to be presented to the user during build
	Description string `json:"description"`
	// Required - indicate if this requires the user to supply a value or not
	Required bool `json:"required"`
	// VerifierRegex - if the user is supplying text and it needs to match a specific pattern, specify a regex pattern here and the UI will indicate to the user if the value is valid or not
	VerifierRegex string `json:"verifier_regex"`
	// DefaultValue - A default value to show the user when building in the Mythic UI. The type here depends on the Parameter Type - ex: for a String, supply a string. For an array, provide an array
	DefaultValue interface{} `json:"default_value"`
	// ParameterType - The type of parameter this is so that the UI can properly render components for the user to modify
	ParameterType BuildParameterType `json:"parameter_type"`
	// FormatString - If Randomize is true, this regex format string is used to generate a value when presenting the option to the user
	FormatString string `json:"format_string"`
	// Randomize - Should this value be randomized each time it's shown to the user so that each payload has a different value
	Randomize bool `json:"randomize"`
	// IsCryptoType -If this is True, then the value supplied by the user is for determining the _kind_ of crypto keys to generate (if any) and the resulting stored value in the database is a dictionary composed of the user's selected and an enc_key and dec_key value
	IsCryptoType bool `json:"crypto_type"`
	// Choices - If the ParameterType is ChooseOne or ChooseMultiple, then the options presented to the user are here.
	Choices []string `json:"choices"`
	// DictionaryChoices - if the ParameterType is Dictionary, then the dictionary choices/preconfigured data is set here
	DictionaryChoices []BuildParameterDictionary `json:"dictionary_choices"`
}

// BuildStep - Identification of a step in the build process that's shown to the user to eventually collect start/end time as well as stdout/stderr per step
type BuildStep struct {
	Name        string `json:"step_name"`
	Description string `json:"step_description"`
}

// PTRPCOtherServiceRPCMessage - A message to call RPC functionality exposed by another Payload Type or C2 Profile
type PTRPCOtherServiceRPCMessage struct {
	// Name - The name of the remote Payload type or C2 Profile
	Name string `json:"service_name"` //required
	// RPCFunction - The name of the function to call for that remote service
	RPCFunction string `json:"service_function"`
	// RPCFunctionArguments - A map of arguments to supply to that remote function
	RPCFunctionArguments map[string]interface{} `json:"service_arguments"`
}

// PTRPCOtherServiceRPCMessageResponse - The result of calling RPC functionality exposed by another Payload Type or C2 Profile
type PTRPCOtherServiceRPCMessageResponse struct {
	// Success - An indicator if the call was successful or not
	Success bool `json:"success"`
	// Error - If the call was unsuccessful, this is an error message about what happened
	Error string `json:"error"`
	// Result - The result returned by the remote service
	Result map[string]interface{} `json:"result"`
}

// PayloadType - The definition of a Payload Type to be synced with Mythic.
/*
	Use the following functions to add an instance of your payload type and build data to Mythic's tracking:
	agentstructs.AllPayloadData.Get("agentname").AddPayloadDefinition(payloadDefinition)
	agentstructs.AllPayloadData.Get("agentname").AddBuildFunction(build)
*/
type PayloadType struct {
	// Name - The name of the payload type that appears in the Mythic UI
	Name string `json:"name"`
	// FileExtension - The default file extension to append to the payload type's name as a placeholder for a filename when generating a payload
	FileExtension string `json:"file_extension"`
	// Author - the name or handle of the author(s) responsible for creating this payload type
	Author string `json:"author"`
	// SupportedOS - An array of operating system names that this payload can compile for. This is used to populate that first dropdown in the Mythic UI when building a payload
	SupportedOS []string `json:"supported_os"`
	// Wrapper - Is this a payload type a wrapper for other payload types or is it a regular payload type
	Wrapper bool `json:"wrapper"`
	// CanBeWrappedByTheFollowingPayloadTypes - Which wrapper payload types does this payload type support (i.e. If this payload type can be supplied to the service_wrapper payload type, list service_wrapper here)
	CanBeWrappedByTheFollowingPayloadTypes []string `json:"supported_wrapper_payload_types"`
	// SupportsDynamicLoading - Does this payload type allow you to dynamically select which commands are loaded into the base payload? If so, set this to True, otherwise all commands are baked into the agent all the time.
	SupportsDynamicLoading bool `json:"supports_dynamic_load"`
	// Description - The description of the payload type to show in the Mythic UI
	Description string `json:"description"`
	// SupportedC2Profiles - The names of the c2 profiles that this payload type supports
	SupportedC2Profiles []string `json:"supported_c2_profiles"`
	// TranslationContainerName - If this payload type uses a translation container, this should be the name of that service
	TranslationContainerName string `json:"translation_container_name"`
	// MythicEncryptsData - If this is True, then Mythic will handle encryption/decryption in messages. If this is false, mythic expects your payload type to have a translation container to handle encryption/decryption on your behalf
	MythicEncryptsData bool `json:"mythic_encrypts"`
	// BuildParameters - A list of build parameters to show to the user during the build process to customize how your payload type's build function operates
	BuildParameters []BuildParameter `json:"build_parameters"`
	// BuildSteps - A list of steps that your build process goes through so that you can report back to the user about the state of the build while it's happening
	BuildSteps []BuildStep `json:"build_steps"`
	// AgentIcon - Don't set this directly, use the agentstructs.AllPayloadData.Get("agentName").AddIcon(filepath.Join(".", "path", "agentname.svg")) call to set this value
	AgentIcon *[]byte `json:"agent_icon"` // automatically filled in based on Name
	// CustomRPCFunctions - The RPC functions you want to expose to other PayloadTypes or C2 Profiles
	CustomRPCFunctions map[string]func(message PTRPCOtherServiceRPCMessage) PTRPCOtherServiceRPCMessageResponse `json:"-"`
	// MessageFormat identifies if the agent uses json or xml messages with Mythic. If you're using a translation container for a custom format, you'd set this to whichever (json/xml) you're going to do your conversions to.
	MessageFormat string `json:"message_format"`
	// AgentType identifies if the payload type is a standard "agent" or if it is another use case like "service" for 3rd party service agents
	AgentType string `json:"agent_type"`
}

// Command - The base definition of a command
/*
	Use the following function to add this command to Mythic's internal tracking:
	agentstructs.AllPayloadData.Get("poseidon").AddCommand(commandDefinition)

	This is easiest to add as part of the init() function for your command file so it's added automatically
*/
type Command struct {
	// Name - the name of the command as the user would type it
	Name string `json:"name"`
	// NeedsAdminPermissions - Does the command need elevated permissions to execute?
	NeedsAdminPermissions bool `json:"needs_admin_permissions"`
	// HelpString - When the user types 'help', what short help would you provide?
	HelpString string `json:"help_string"`
	// Description - A description of what the command does that appears in the tasking modal as well as when the user is selecting commands to include in their payload
	Description string `json:"description"`
	// Version - What version of this command is this? The version is tracked overall and per-load within a Payload and Callback. This makes it easier to see if a callback or payload has an outdated version of a command.
	Version uint32 `json:"version"`
	// SupportedUIFeatures - The list of UI features that the command supports such as 'callback_table:exit` or `file_browser:list`.
	/*
		The most common of these features can be found on the Mythic documentation website, but you can make your own custom ones as well.
		When you want to do browser scripting and support issuing a task with a button click, that task is identified based on the supported_ui_features you supply here.
		There's no required format, but typically they're in the form of `general:specific`, so maybe `registry:write` or `clipboard:set`.
	*/
	SupportedUIFeatures []string `json:"supported_ui_features"`
	// Author - the author(s) of this command
	Author string `json:"author"`
	// MitreAttackMappings - A list of MITRE Technique IDs (ex: T1033) that this command maps to
	MitreAttackMappings []string `json:"attack"`
	// ScriptOnlyCommand - Is this command only defined as a script/golang file or does it have a matching function within the payload
	ScriptOnlyCommand bool `json:"script_only"`
	// CommandAttributes - Attributes about this command that can be used to determine what commands the user can select when building the payload.
	// This also comes into play when determining commands to list for some command parameters
	CommandAttributes CommandAttribute `json:"attributes"`
	// CommandParameters - A list of the parameters associated with this command (also known as arguments)
	CommandParameters []CommandParameter `json:"parameters"`
	// AssociatedBrowserScript - If this command has a browser script to manipulate the output from this command, reference that here
	AssociatedBrowserScript *BrowserScript `json:"browserscript,omitempty"`
	// TaskFunctionOPSECPre - If you want to provide an OPSEC check before your main TaskFunctionCreateTasking function, you can define that function here
	TaskFunctionOPSECPre PtTaskFunctionOPSECPre `json:"-"`
	// TaskFunctionCreateTasking - This is the main function to do additional processing, RPC calls, and anything else before your command is ready for the agent to pick it up
	TaskFunctionCreateTasking PtTaskFunctionCreateTasking `json:"-"`
	// TaskFunctionProcessResponse - If your callback returns data in the 'process_response' key within your responses array, that data gets processed here.
	TaskFunctionProcessResponse PtTaskFunctionProcessResponse `json:"-"`
	// TaskFunctionOPSECPost - If you want to provide an OPSEC check after your TaskFunctionCreateTasking function executes but before the agent picks up the tasking, you can do that here
	TaskFunctionOPSECPost PtTaskFunctionOPSECPost `json:"-"`
	// TaskFunctionParseArgString - Parse an argument string from the user into your command's CommandParameters array
	TaskFunctionParseArgString PtTaskFunctionParseArgString `json:"-"`
	// TaskFunctionParseArgDictionary - Parse an argument dictionary from the user into your command's CommandParameters array
	TaskFunctionParseArgDictionary PtTaskFunctionParseArgDictionary `json:"-"`
	// TaskCompletionFunctions - If your TaskFunctionCreateTasking function or any of your subtasks have completion functions, define them here
	TaskCompletionFunctions map[string]PTTaskCompletionFunction `json:"-"`
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
	COMMAND_PARAMETER_TYPE_TYPED_ARRAY                          = "TypedArray"
)

// CommandParameter - The base definition for a parameter (i.e. argument) to your command
type CommandParameter struct {
	// Name - The name of your parameter - used when adding args or changing arg values
	Name string `json:"name"`
	// ModalDisplayName - A more friendly version of the name, most likely with captialization and spaces
	ModalDisplayName string `json:"display_name"`
	// CLIName - A more CLI friendly version of the name, potentially without dashes/underscores and no spaces
	CLIName string `json:"cli_name"`
	// ParameterType - The type of parameter - this influences how things work in the UI
	ParameterType CommandParameterType `json:"parameter_type"`
	// Description - The description of the parameter that's displayed to the user when they hover over the ModalDisplayName
	Description string `json:"description"`
	// Choices - If the ParameterType is ChooseOne or ChooseMultiple, these are the choices for the user.
	// If the ParameterType is TypedArray, these are the options for each array entry
	Choices []string `json:"choices"`
	// DefaultValue - The default value to present to the user when they pull up the modal view
	DefaultValue interface{} `json:"default_value"`
	// SupportedAgents - When using the "Payload" Parameter Type, you can filter down which payloads are presented to the operator based on this list of supported agents.
	SupportedAgents []string `json:"supported_agents"`
	// SupportedAgentBuildParameters - When using the "Payload" Parameter Type, you can filter down which payloads are presented to the operator based on specific build parameters for specific payload types.
	SupportedAgentBuildParameters map[string]string `json:"supported_agent_build_parameters"`
	// ChoicesAreAllCommands - Can be used with ChooseOne or ChooseMultiple Parameter Types to automatically populate those options in the UI with all of the commands for the payload type.
	ChoicesAreAllCommands bool `json:"choices_are_all_commands"`
	// ChoicesAreLoadedCommands - Can be used with ChooseOne or ChooseMultiple Parameter Types to automatically populate those options in the UI with all of the currently loaded commands.
	ChoicesAreLoadedCommands bool `json:"choices_are_loaded_commands"`
	// FilterCommandChoicesByCommandAttributes -  When using the ChooseOne or ChooseMultiple Parameter type along with choices_are_all_commands, you can filter down those options based on attribute values in your command's CommandAttributes field.
	FilterCommandChoicesByCommandAttributes map[string]string `json:"choice_filter_by_command_attributes"`
	// DynamicQueryFunction -  Provide a dynamic query function to be called when the user views that parameter option in the UI to populate choices for the ChooseOne or ChooseMultiple Parameter Types.
	DynamicQueryFunction PTTaskingDynamicQueryFunction `json:"dynamic_query_function"`
	// TypedArrayParseFunction - Provide a function to be called when the user types out a typedArray value on the CLI, but that needs to be parsed for a Modal Popup
	TypedArrayParseFunction PTTaskingTypedArrayParseFunction `json:"typedarray_parse_function"`
	// ParameterGroupInformation - Define 0+ different parameter groups that this parameter belongs to.
	ParameterGroupInformation []ParameterGroupInfo `json:"parameter_group_info"`
	value                     interface{}          // the current value for the parameter
	userSupplied              bool                 // was this value supplied by the user or a default value
}

type PTTaskingDynamicQueryFunction func(PTRPCDynamicQueryFunctionMessage) []string
type PTTaskingTypedArrayParseFunction func(message PTRPCTypedArrayParseFunctionMessage) [][]string

func (f PTTaskingDynamicQueryFunction) MarshalJSON() ([]byte, error) {
	if f != nil {
		return json.Marshal("function defined")
	} else {
		return json.Marshal("")
	}
}
func (f PTTaskingTypedArrayParseFunction) MarshalJSON() ([]byte, error) {
	if f != nil {
		return json.Marshal("function defined")
	} else {
		return json.Marshal("")
	}
}

// CommandAttribute - Attributes about a specific command to influence build options and command parameter options
type CommandAttribute struct {
	// SupportedOS -  Which operating systems does this command support? An empty list means all OS.
	SupportedOS []string `json:"supported_os"`
	// CommandIsBuiltin -  Is this command baked into the agent permanently?
	CommandIsBuiltin bool `json:"builtin"`
	// CommandIsSuggested - If true, this command will appear on the "included" side when building your payload by default.
	CommandIsSuggested bool `json:"suggested_command"`
	// CommandCanOnlyBeLoadedLater - If true, this command can only be loaded after you have a callback and not included in the base payload.
	CommandCanOnlyBeLoadedLater bool `json:"load_only"`
	// FilterCommandAvailabilityByAgentBuildParameters - Specify if this command is allowed to be built into the payload or not based on build parameters the user specifies.
	/*
		is of the form {"build param name": "build param value"}
	*/
	FilterCommandAvailabilityByAgentBuildParameters map[string]string `json:"filter_by_build_parameter"`
	// AdditionalAttributes - Additional, developer-supplied, key-value pairs such as a dependency note that a command relies on another comand
	AdditionalAttributes map[string]string `json:"additional_items"`
}

// ParameterGroupInfo - Allow conditional parameters displayed to the user and determine if parameters are required and the order in which they're presented to the user
type ParameterGroupInfo struct {
	// ParameterIsRequired - Is this parameter required?
	ParameterIsRequired bool `json:"required"`
	// GroupName - What is the name of this parameter group (i.e. group of parameters that are grouped together)
	GroupName string `json:"group_name"`
	// UIModalPosition - If the user opens a modal to fill out parameters, which position should this parameter be shown?
	UIModalPosition uint32 `json:"ui_position"`
	// AdditionalInformation - Additional, developer-supplied, key-value pairs of information
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
	Description        string                                `json:"description"`
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
	IsP2P      bool                   `json:"c2_profile_is_p2p"`
	Parameters map[string]interface{} `json:"c2_profile_parameters"`
}
type PayloadConfigurationBuildParameter struct {
	Name  string      `json:"name" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

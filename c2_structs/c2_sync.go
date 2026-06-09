package c2structs

import (
	"context"
	"encoding/json"

	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

// C2_SYNC STRUCTS
type C2ParameterType = string

type C2SyncMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type C2SyncMessage struct {
	Profile          C2Profile     `json:"c2_profile"`
	Parameters       []C2Parameter `json:"parameters"`
	ContainerVersion string        `json:"container_version"`
}
type C2RPCOtherServiceRPCMessage struct {
	ServiceName                 string                 `json:"service_name"` //required
	ServiceRPCFunction          string                 `json:"service_function"`
	ServiceRPCFunctionArguments map[string]interface{} `json:"service_arguments"`
}
type C2RPCOtherServiceRPCMessageResponse struct {
	Success               bool                   `json:"success"`
	Error                 string                 `json:"error"`
	Result                map[string]interface{} `json:"result"`
	RestartInternalServer bool                   `json:"restart_internal_server"`
}

type C2Profile struct {
	Name                       string                                                                                                     `json:"name"`
	Description                string                                                                                                     `json:"description"`
	Author                     string                                                                                                     `json:"author"`
	IsP2p                      bool                                                                                                       `json:"is_p2p"`
	IsServerRouted             bool                                                                                                       `json:"is_server_routed"`
	ServerBinaryPath           string                                                                                                     `json:"-"`
	ServerFolderPath           string                                                                                                     `json:"-"`
	SemVer                     string                                                                                                     `json:"semver"`
	AgentIcon                  *[]byte                                                                                                    `json:"agent_icon"`
	DarkModeAgentIcon          *[]byte                                                                                                    `json:"dark_mode_agent_icon"`
	ConfigCheckFunction        func(context.Context, C2ConfigCheckMessage) C2ConfigCheckMessageResponse                                   `json:"-"`
	GetRedirectorRulesFunction func(context.Context, C2GetRedirectorRuleMessage) C2GetRedirectorRuleMessageResponse                       `json:"-"`
	OPSECCheckFunction         func(context.Context, C2OPSECMessage) C2OPSECMessageResponse                                               `json:"-"`
	GetIOCFunction             func(context.Context, C2GetIOCMessage) C2GetIOCMessageResponse                                             `json:"-"`
	SampleMessageFunction      func(context.Context, C2SampleMessageMessage) C2SampleMessageResponse                                      `json:"-"`
	HostFileFunction           func(context.Context, C2HostFilesMessage) C2HostFilesMessageResponse                                       `json:"-"`
	CustomRPCFunctions         map[string]func(context.Context, C2RPCOtherServiceRPCMessage) C2RPCOtherServiceRPCMessageResponse          `json:"-"`
	OnContainerStartFunction   func(context.Context, sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse `json:"-"`
}

const (
	C2_PARAMETER_TYPE_STRING            C2ParameterType = "String"
	C2_PARAMETER_TYPE_BOOLEAN                           = "Boolean"
	C2_PARAMETER_TYPE_CHOOSE_ONE                        = "ChooseOne"
	C2_PARAMETER_TYPE_CHOOSE_ONE_CUSTOM                 = "ChooseOneCustom"
	C2_PARAMETER_TYPE_CHOOSE_MULTIPLE                   = "ChooseMultiple"
	C2_PARAMETER_TYPE_ARRAY                             = "Array"
	C2_PARAMETER_TYPE_DATE                              = "Date"
	C2_PARAMETER_TYPE_DICTIONARY                        = "Dictionary"
	C2_PARAMETER_TYPE_NUMBER                            = "Number"
	C2_PARAMETER_TYPE_FILE                              = "File"
	C2_PARAMETER_TYPE_FILE_MULTIPLE                     = "FileMultiple"
	C2_PARAMETER_TYPE_JSON_STRING                       = "JSONString"
)

type HideConditionOperand string

const (
	HideConditionOperandEQ                 HideConditionOperand = "eq"
	HideConditionOperandNotEQ                                   = "neq"
	HideConditionOperandIN                                      = "in"
	HideConditionOperandNotIN                                   = "nin"
	HideConditionOperandLessThan                                = "lt"
	HideConditionOperandGreaterThan                             = "gt"
	HideConditionOperandLessThanOrEqual                         = "lte"
	HideConditionOperandGreaterThanOrEqual                      = "gte"
	HideConditionOperationStartsWith                            = "sw"
	HideConditionOperationEndsWith                              = "ew"
	HideConditionOperationContains                              = "co"
	HideConditionOperationNotContains                           = "nco"
)

type BuildParameterHideCondition struct {
	Name    string               `json:"name"`
	Operand HideConditionOperand `json:"operand"`
	Value   string               `json:"value"`
	Choices []string             `json:"choices"`
}
type C2RPCC2ParameterDynamicQueryFunction func(context.Context, C2RPCDynamicQueryC2ParameterFunctionMessage) C2RPCDynamicQueryC2ParameterFunctionMessageResponse

func (f C2RPCC2ParameterDynamicQueryFunction) MarshalJSON() ([]byte, error) {
	if f != nil {
		return json.Marshal("function defined")
	}
	return json.Marshal("")
}

type C2Parameter struct {
	Name                 string                               `json:"name"`
	DisplayName          string                               `json:"display_name"`
	Description          string                               `json:"description"`
	Required             bool                                 `json:"required"`
	VerifierRegex        string                               `json:"verifier_regex"`
	DefaultValue         interface{}                          `json:"default_value"`
	ParameterType        C2ParameterType                      `json:"parameter_type"`
	FormatString         string                               `json:"format_string"`
	Randomize            bool                                 `json:"randomize"`
	IsCryptoType         bool                                 `json:"crypto_type"`
	Choices              []string                             `json:"choices"`
	ChoicesDisplayNames  map[string]string                    `json:"choices_display_names"`
	DictionaryChoices    []C2ParameterDictionary              `json:"dictionary_choices"`
	JsonStringSchema     map[string]interface{}               `json:"json_string_schema"`
	GroupName            string                               `json:"group_name"`
	HideConditions       []BuildParameterHideCondition        `json:"hide_conditions"`
	UiPosition           int                                  `json:"ui_position"`
	DynamicQueryFunction C2RPCC2ParameterDynamicQueryFunction `json:"dynamic_query_function"`
}

type C2ParameterDictionary struct {
	Name         string `json:"name"`
	DefaultValue string `json:"default_value"`
	DefaultShow  bool   `json:"default_show"`
}

type SimplifiedC2ParameterDictionary struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Key   string `json:"key"`
}

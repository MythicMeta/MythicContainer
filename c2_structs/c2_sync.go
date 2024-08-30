package c2structs

import "github.com/MythicMeta/MythicContainer/utils/sharedStructs"

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
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Result  map[string]interface{} `json:"result"`
}

type C2Profile struct {
	Name                       string                                                                                    `json:"name"`
	Description                string                                                                                    `json:"description"`
	Author                     string                                                                                    `json:"author"`
	IsP2p                      bool                                                                                      `json:"is_p2p"`
	IsServerRouted             bool                                                                                      `json:"is_server_routed"`
	ServerBinaryPath           string                                                                                    `json:"-"`
	ServerFolderPath           string                                                                                    `json:"-"`
	ConfigCheckFunction        func(message C2ConfigCheckMessage) C2ConfigCheckMessageResponse                           `json:"-"`
	GetRedirectorRulesFunction func(message C2GetRedirectorRuleMessage) C2GetRedirectorRuleMessageResponse               `json:"-"`
	OPSECCheckFunction         func(message C2OPSECMessage) C2OPSECMessageResponse                                       `json:"-"`
	GetIOCFunction             func(message C2GetIOCMessage) C2GetIOCMessageResponse                                     `json:"-"`
	SampleMessageFunction      func(message C2SampleMessageMessage) C2SampleMessageResponse                              `json:"-"`
	HostFileFunction           func(message C2HostFileMessage) C2HostFileMessageResponse                                 `json:"-"`
	CustomRPCFunctions         map[string]func(message C2RPCOtherServiceRPCMessage) C2RPCOtherServiceRPCMessageResponse  `json:"-"`
	OnContainerStartFunction   func(sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse `json:"-"`
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
	C2_PARAMETER_TYPE_TYPED_ARRAY                       = "TypedArray"
	C2_PARAMETER_TYPE_FILE                              = "File"
	C2_PARAMETER_TYPE_FILE_MULTIPLE                     = "FileMultiple"
)

type C2Parameter struct {
	Description       string                  `json:"description"`
	Name              string                  `json:"name"`
	DefaultValue      interface{}             `json:"default_value"`
	Randomize         bool                    `json:"randomize"`
	FormatString      string                  `json:"format_string"`
	ParameterType     C2ParameterType         `json:"parameter_type"`
	Required          bool                    `json:"required"`
	VerifierRegex     string                  `json:"verifier_regex"`
	IsCryptoType      bool                    `json:"crypto_type"`
	Choices           []string                `json:"choices"`
	DictionaryChoices []C2ParameterDictionary `json:"dictionary_choices"`
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

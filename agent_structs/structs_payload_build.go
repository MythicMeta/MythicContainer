package agentstructs

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)
import "errors"

// PAYLOAD_BUILD STRUCTS

// PayloadBuildMessage - A structure of the build information the user provided to generate an instance of the payload type.
// This information gets passed to your payload type's build function.
type PayloadBuildMessage struct {
	// PayloadType - the name of the payload type for the build
	PayloadType string `json:"payload_type" mapstructure:"payload_type"`
	// Filename - the name of the file the user originally supplied for this build
	Filename string `json:"filename" mapstructure:"filename"`
	// CommandList - the list of commands the user selected to include in the build
	CommandList []string `json:"commands" mapstructure:"commands"`
	// build param name : build value
	// BuildParameters - map of param name -> build value from the user for the build parameters defined
	// File type build parameters are supplied as a string UUID to use with MythicRPC for fetching file contents
	// Array type build parameters are supplied as []string{}
	BuildParameters
	// C2Profiles - list of C2 profiles selected to include in the payload and their associated parameters
	C2Profiles []PayloadBuildC2Profile `json:"c2profiles" mapstructure:"c2profiles"`
	// WrappedPayload - bytes of the wrapped payload if one exists
	WrappedPayload *[]byte `json:"wrapped_payload,omitempty" mapstructure:"wrapped_payload"`
	// WrappedPayloadUUID - the UUID of the wrapped payload if one exists
	WrappedPayloadUUID *string `json:"wrapped_payload_uuid,omitempty" mapstructure:"wrapped_payload_uuid"`
	// SelectedOS - the operating system the user selected when building the agent
	SelectedOS string `json:"selected_os" mapstructure:"selected_os"`
	// PayloadUUID - the Mythic generated UUID for this payload instance
	PayloadUUID string `json:"uuid" mapstructure:"uuid"`
	// PayloadFileUUID - The Mythic generated File UUID associated with this payload
	PayloadFileUUID string `json:"payload_file_uuid" mapstructure:"payload_file_uuid"`
	// Secrets - User supplied secrets that get sent down with payload builds
	Secrets map[string]interface{} `json:"secrets"`
}

// PayloadBuildC2Profile - A structure of the selected C2 Profile information the user selected to build into a payload.
type PayloadBuildC2Profile struct {
	Name  string `json:"name" mapstructure:"name"`
	IsP2P bool   `json:"is_p2p" mapstructure:"is_p2p"`
	// parameter name: parameter value
	// Parameters - this is an interface of parameter name -> parameter value from the associated C2 profile.
	// The types for the various parameter names can be found by looking at the build parameters in the Mythic UI.
	Parameters map[string]interface{} `json:"parameters" mapstructure:"parameters"`
}

type CryptoArg struct {
	Value  string `json:"value" mapstructure:"value"`
	EncKey string `json:"enc_key" mapstructure:"enc_key"`
	DecKey string `json:"dec_key" mapstructure:"dec_key"`
}

func (arg *PayloadBuildC2Profile) GetArg(name string) (interface{}, error) {
	for key, currentArg := range arg.Parameters {
		if key == name {
			return currentArg, nil
		}
	}
	return nil, errors.New("failed to find argument")
}
func (arg *PayloadBuildC2Profile) GetArgNames() []string {
	argNames := []string{}
	for key, _ := range arg.Parameters {
		argNames = append(argNames, key)
	}
	return argNames
}
func (arg *PayloadBuildC2Profile) GetStringArg(name string) (string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return "", err
	} else if val == nil {
		return "", nil
	} else {
		return getTypedValue[string](val)
	}
}
func (arg *PayloadBuildC2Profile) GetNumberArg(name string) (float64, error) {
	if val, err := arg.GetArg(name); err != nil {
		return 0, err
	} else if val == nil {
		return 0, nil
	} else if floatVal, err := getTypedValue[float64](val); err == nil {
		return floatVal, nil
	} else if intVal, err := getTypedValue[int](val); err == nil {
		return float64(intVal), nil
	} else {
		return 0, err
	}
}
func (arg *PayloadBuildC2Profile) GetBooleanArg(name string) (bool, error) {
	if val, err := arg.GetArg(name); err != nil {
		return false, err
	} else if val == nil {
		return false, nil
	} else {
		return getTypedValue[bool](val)
	}
}
func (arg *PayloadBuildC2Profile) GetDictionaryArg(name string) (map[string]string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return nil, err
	} else if val == nil {
		return make(map[string]string), nil
	} else if initialDict, err := getTypedValue[map[string]interface{}](val); err != nil {
		return nil, err
	} else {
		finalMap := make(map[string]string, len(initialDict))
		for key, val := range initialDict {
			switch v := val.(type) {
			case string:
				finalMap[key] = v
			default:
				finalMap[key] = fmt.Sprintf("%v", v)
			}
		}
		return finalMap, nil
	}
}
func (arg *PayloadBuildC2Profile) GetChooseOneArg(name string) (string, error) {
	return arg.GetStringArg(name)
}
func (arg *PayloadBuildC2Profile) GetArrayArg(name string) ([]string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return []string{}, err
	} else if val == nil {
		return []string{}, nil
	} else if interfaceArray, err := getTypedValue[[]interface{}](val); err != nil {
		return []string{}, err
	} else {
		stringArray := make([]string, len(interfaceArray))
		for index, _ := range interfaceArray {
			stringArray[index] = fmt.Sprintf("%v", interfaceArray[index])
		}
		return stringArray, nil
	}
}
func (arg *PayloadBuildC2Profile) GetChooseMultipleArg(name string) ([]string, error) {
	return arg.GetArrayArg(name)
}
func (arg *PayloadBuildC2Profile) GetDateArg(name string) (string, error) {
	return arg.GetStringArg(name)
}
func (arg *PayloadBuildC2Profile) GetFileArg(name string) (string, error) {
	return arg.GetStringArg(name)
}
func (arg *PayloadBuildC2Profile) GetCryptoArg(name string) (CryptoArg, error) {
	cryptoArg := CryptoArg{}
	if val, err := arg.GetArg(name); err != nil {
		return cryptoArg, err
	} else if val == nil {
		return cryptoArg, nil
	} else if err := mapstructure.Decode(val, &cryptoArg); err != nil {
		return cryptoArg, err
	} else {
		return cryptoArg, nil
	}
}
func (arg *PayloadBuildC2Profile) GetTypedArrayArg(name string) ([][]string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return [][]string{}, err
	} else if val == nil {
		return [][]string{}, nil
	} else if interfaceArray, err := getTypedValue[[][]interface{}](val); err != nil {
		return [][]string{}, err
	} else {
		stringArray := make([][]string, len(interfaceArray))
		for index, _ := range interfaceArray {
			stringArray[index] = []string{}
			for index2, _ := range interfaceArray[index] {
				stringArray[index] = append(stringArray[index], fmt.Sprintf("%v", interfaceArray[index][index2]))
			}
		}
		return stringArray, nil
	}
}

type BuildParameters struct {
	Parameters map[string]interface{} `json:"build_parameters" mapstructure:"build_parameters"`
}

func (arg *BuildParameters) GetArg(name string) (interface{}, error) {
	for key, currentArg := range arg.Parameters {
		if key == name {
			return currentArg, nil
		}
	}
	return nil, errors.New("failed to find argument")
}
func (arg *BuildParameters) GetArgNames() []string {
	argNames := []string{}
	for key, _ := range arg.Parameters {
		argNames = append(argNames, key)
	}
	return argNames
}
func (arg *BuildParameters) GetStringArg(name string) (string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return "", err
	} else if val == nil {
		return "", nil
	} else {
		return getTypedValue[string](val)
	}
}
func (arg *BuildParameters) GetNumberArg(name string) (float64, error) {
	if val, err := arg.GetArg(name); err != nil {
		return 0, err
	} else if val == nil {
		return 0, nil
	} else if floatVal, err := getTypedValue[float64](val); err == nil {
		return floatVal, nil
	} else if intVal, err := getTypedValue[int](val); err == nil {
		return float64(intVal), nil
	} else {
		return 0, err
	}
}
func (arg *BuildParameters) GetBooleanArg(name string) (bool, error) {
	if val, err := arg.GetArg(name); err != nil {
		return false, err
	} else if val == nil {
		return false, nil
	} else {
		return getTypedValue[bool](val)
	}
}
func (arg *BuildParameters) GetDictionaryArg(name string) (map[string]string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return nil, err
	} else if val == nil {
		return make(map[string]string), nil
	} else if initialDict, err := getTypedValue[map[string]interface{}](val); err != nil {
		return nil, err
	} else {
		finalMap := make(map[string]string, len(initialDict))
		for key, val := range initialDict {
			switch v := val.(type) {
			case string:
				finalMap[key] = v
			default:
				finalMap[key] = fmt.Sprintf("%v", v)
			}
		}
		return finalMap, nil
	}
}
func (arg *BuildParameters) GetChooseOneArg(name string) (string, error) {
	return arg.GetStringArg(name)
}
func (arg *BuildParameters) GetArrayArg(name string) ([]string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return []string{}, err
	} else if val == nil {
		return []string{}, nil
	} else if interfaceArray, err := getTypedValue[[]interface{}](val); err != nil {
		return []string{}, err
	} else {
		stringArray := make([]string, len(interfaceArray))
		for index, _ := range interfaceArray {
			stringArray[index] = fmt.Sprintf("%v", interfaceArray[index])
		}
		return stringArray, nil
	}
}
func (arg *BuildParameters) GetChooseMultipleArg(name string) ([]string, error) {
	return arg.GetArrayArg(name)
}
func (arg *BuildParameters) GetDateArg(name string) (string, error) {
	return arg.GetStringArg(name)
}
func (arg *BuildParameters) GetFileArg(name string) (string, error) {
	return arg.GetStringArg(name)
}
func (arg *BuildParameters) GetCryptoArg(name string) (CryptoArg, error) {
	cryptoArg := CryptoArg{}
	if val, err := arg.GetArg(name); err != nil {
		return cryptoArg, err
	} else if val == nil {
		return cryptoArg, nil
	} else if err := mapstructure.Decode(val, &cryptoArg); err != nil {
		return cryptoArg, err
	} else {
		return cryptoArg, nil
	}
}
func (arg *BuildParameters) GetTypedArrayArg(name string) ([][]string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return [][]string{}, err
	} else if val == nil {
		return [][]string{}, nil
	} else if interfaceArray, err := getTypedValue[[][]interface{}](val); err != nil {
		return [][]string{}, err
	} else {
		stringArray := make([][]string, len(interfaceArray))
		for index, _ := range interfaceArray {
			stringArray[index] = []string{}
			for index2, _ := range interfaceArray[index] {
				stringArray[index] = append(stringArray[index], fmt.Sprintf("%v", interfaceArray[index][index2]))
			}
		}
		return stringArray, nil
	}
}

type PAYLOAD_BUILD_STATUS = string

const (
	PAYLOAD_BUILD_STATUS_SUCCESS PAYLOAD_BUILD_STATUS = "success"
	PAYLOAD_BUILD_STATUS_ERROR                        = "error"
)

// PayloadBuildResponse - The result of calling a payload type's build function. This returns not only the actual
// payload bytes, but surrounding metadata such as updated filenames, command lists, and stdout/stderr messages.
type PayloadBuildResponse struct {
	// PayloadUUID - The UUID associated with this payload
	PayloadUUID string `json:"uuid"`
	// Success - was this build process successful or not
	Success bool `json:"success"`
	// UpdatedFilename - Optionally updated filename based on build parameters to more closely match the return file type
	UpdatedFilename *string `json:"updated_filename,omitempty"`
	// Payload - the raw bytes of the payload that was compiled/created
	Payload *[]byte `json:"payload,omitempty"`
	// UpdatedCommandList - if you want to adjust the list of commands in this payload from what the user provided,
	// provide the updated list of command names here
	UpdatedCommandList *[]string `json:"updated_command_list,omitempty"`
	// BuildStdErr - build stderr message to associate with the build
	BuildStdErr string `json:"build_stderr"`
	// BuildStdOut - build stdout message to associate with the build
	BuildStdOut string `json:"build_stdout"`
	// BuildMessage - general message to associate with the build. Usually not as verbose as the stdout/stderr.
	BuildMessage string `json:"build_message"`
}

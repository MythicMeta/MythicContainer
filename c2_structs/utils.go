package c2structs

import (
	"errors"
	"fmt"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"github.com/mitchellh/mapstructure"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type RabbitmqRPCMethod struct {
	RabbitmqRoutingKey         string
	RabbitmqProcessingFunction func([]byte) interface{}
}
type RabbitmqDirectMethod struct {
	RabbitmqRoutingKey         string
	RabbitmqProcessingFunction func([]byte)
}
type allC2Data struct {
	c2Definition         C2Profile
	parameters           []C2Parameter
	mutex                sync.RWMutex
	containerVersion     string
	rpcMethods           []RabbitmqRPCMethod
	directMethods        []RabbitmqDirectMethod
	RunningServerProcess *exec.Cmd
	OutputChannel        chan string
}

var (
	AllC2Data containerC2Data
)

type containerC2Data struct {
	C2Map map[string]*allC2Data
}

func (r *containerC2Data) GetAllNames() []string {
	names := []string{}
	for key, _ := range r.C2Map {
		if key != "" && !helpers.StringSliceContains(names, key) {
			names = append(names, key)
		}

	}
	return names
}
func (r *containerC2Data) Get(name string) *allC2Data {
	if r.C2Map == nil {
		r.C2Map = make(map[string]*allC2Data)
	}
	if existingC2Data, ok := r.C2Map[name]; !ok {
		newC2Data := allC2Data{}
		newC2Data.OutputChannel = make(chan string, 20)
		r.C2Map[name] = &newC2Data
		return &newC2Data
	} else {
		return existingC2Data
	}
}
func (r *allC2Data) AddC2Definition(payloadDef C2Profile) {

	if payloadDef.ServerFolderPath == "" {
		if payloadDef.ServerBinaryPath == "" {
			if osPath, err := os.Executable(); err != nil {
				logging.LogError(err, "Failed to get the current working path")
				os.Exit(1)
			} else {
				payloadDef.ServerFolderPath = filepath.Dir(osPath)
			}
		} else if serverFilePath, err := filepath.Abs(payloadDef.ServerBinaryPath); err != nil {
			logging.LogError(err, "Failed to get the absolute path for the server binary")
		} else {
			payloadDef.ServerFolderPath = filepath.Dir(serverFilePath)
		}
	}

	if payloadDef.IsP2p {
		if payloadDef.ServerFolderPath == "" {
			logging.LogError(nil, "Must supply ServerFolderPath ")
			os.Exit(1)
		}
	} else if payloadDef.ServerBinaryPath == "" {
		logging.LogError(nil, "Must supply ServerBinaryPath")
		os.Exit(1)
	} else if payloadDef.ServerFolderPath == "" {
		logging.LogError(nil, "Failed to get ServerFolderPath from ServerBinaryPath")
		os.Exit(1)
	}
	if payloadDef.CustomRPCFunctions == nil {
		payloadDef.CustomRPCFunctions = make(map[string]func(message C2RPCOtherServiceRPCMessage) C2RPCOtherServiceRPCMessageResponse)
	}
	r.c2Definition = payloadDef
}
func (r *allC2Data) GetC2Definition() C2Profile {
	return r.c2Definition
}
func (r *allC2Data) GetC2ServerFolderPath() string {
	return r.c2Definition.ServerFolderPath
}
func (r *allC2Data) AddParameters(params []C2Parameter) {
	r.parameters = params
}
func (r *allC2Data) GetParameters() []C2Parameter {
	return r.parameters
}
func (r *allC2Data) AddContainerVersion(ver string) {
	r.containerVersion = ver
}
func (r *allC2Data) GetC2Name() string {
	return r.c2Definition.Name
}
func (r *allC2Data) GetContainerVersion() string {
	return r.containerVersion
}
func (r *allC2Data) GetServerPath() string {
	logging.LogInfo("getting server binary path", "path", r.c2Definition.ServerBinaryPath)
	return r.c2Definition.ServerBinaryPath
}
func (r *allC2Data) AddRPCMethod(m RabbitmqRPCMethod) {
	r.mutex.Lock()
	r.rpcMethods = append(r.rpcMethods, m)
	r.mutex.Unlock()
}
func (r *allC2Data) GetRPCMethods() []RabbitmqRPCMethod {
	return r.rpcMethods
}
func (r *allC2Data) AddDirectMethod(m RabbitmqDirectMethod) {
	r.mutex.Lock()
	r.directMethods = append(r.directMethods, m)
	r.mutex.Unlock()
}
func (r *allC2Data) GetDirectMethods() []RabbitmqDirectMethod {
	return r.directMethods
}
func (r *allC2Data) GetRoutingKey(baseRoutingKey string) string {
	return fmt.Sprintf("%s_%s", r.GetC2Name(), baseRoutingKey)
}

type CryptoArg struct {
	Value  string `json:"value" mapstructure:"value"`
	EncKey string `json:"enc_key" mapstructure:"enc_key"`
	DecKey string `json:"dec_key" mapstructure:"dec_key"`
}

type C2Parameters struct {
	Name       string                 `json:"c2_profile_name"`
	Parameters map[string]interface{} `json:"parameters"`
}

func (arg *C2Parameters) GetArg(name string) (interface{}, error) {
	for key, currentArg := range arg.Parameters {
		if key == name {
			return currentArg, nil
		}
	}
	return nil, errors.New("failed to find argument")
}
func (arg *C2Parameters) GetArgNames() []string {
	argNames := []string{}
	for key, _ := range arg.Parameters {
		argNames = append(argNames, key)
	}
	return argNames
}
func (arg *C2Parameters) GetStringArg(name string) (string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return "", err
	} else if val == nil {
		return "", nil
	} else {
		return getTypedValue[string](val)
	}
}
func (arg *C2Parameters) GetNumberArg(name string) (float64, error) {
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
func (arg *C2Parameters) GetBooleanArg(name string) (bool, error) {
	if val, err := arg.GetArg(name); err != nil {
		return false, err
	} else if val == nil {
		return false, nil
	} else {
		return getTypedValue[bool](val)
	}
}
func (arg *C2Parameters) GetDictionaryArg(name string) (map[string]string, error) {
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
func (arg *C2Parameters) GetChooseOneArg(name string) (string, error) {
	return arg.GetStringArg(name)
}
func (arg *C2Parameters) GetArrayArg(name string) ([]string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return []string{}, err
	} else if val == nil {
		return []string{}, nil
	} else {
		return getTypedValue[[]string](val)
	}
}
func (arg *C2Parameters) GetChooseMultipleArg(name string) ([]string, error) {
	return arg.GetArrayArg(name)
}
func (arg *C2Parameters) GetDateArg(name string) (string, error) {
	return arg.GetStringArg(name)
}
func (arg *C2Parameters) GetFileArg(name string) (string, error) {
	return arg.GetStringArg(name)
}
func (arg *C2Parameters) GetCryptoArg(name string) (CryptoArg, error) {
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
func getTypedValue[T any](value interface{}) (T, error) {
	switch v := value.(type) {
	case T:
		return v, nil
	default:
		var emptyResult T
		return emptyResult, errors.New("bad type for value")
	}
}

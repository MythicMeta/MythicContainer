package agentstructs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// Args helper functions
func GenerateArgsData(cmdParams []CommandParameter, task PTTaskMessageAllData) (PTTaskMessageArgsData, error) {
	args := PTTaskMessageArgsData{
		taskingLocation:           task.Task.TaskingLocation,
		commandLine:               task.Task.Params,
		rawCommandLine:            task.Task.OriginalParams,
		initialParameterGroupName: task.Task.ParameterGroupName,
	}
	//fmt.Printf("parameter group name: %s\n", task.Task.ParameterGroupName)
	for paramIndex := range cmdParams {
		param := cmdParams[paramIndex]
		for groupInfoIndex := range param.ParameterGroupInformation {
			if param.ParameterGroupInformation[groupInfoIndex].GroupName == "" {
				param.ParameterGroupInformation[groupInfoIndex].GroupName = "Default"
			}
		}
		args.args = append(args.args, param)
	}
	//logging.LogInfo("just called GenerateArgsData", "args", args.args, "cmdParams", cmdParams, "command", task.Task.OriginalParams)
	return args, nil
}
func (cmd *CommandParameter) GetCurrentValue() interface{} {
	if cmd.value == nil {
		return cmd.DefaultValue
	} else {
		return cmd.value
	}
}
func (arg *PTTaskMessageArgsData) SetManualArgs(args string) {
	arg.manualArgs = &args
}
func (arg *PTTaskMessageArgsData) SetManualParameterGroup(groupName string) {
	arg.manualParameterGroupName = groupName
}
func (arg *PTTaskMessageArgsData) GetArg(name string) (interface{}, error) {
	for _, currentArg := range arg.args {
		if currentArg.Name == name || currentArg.CLIName == name {
			return currentArg.GetCurrentValue(), nil
		}
	}
	return nil, errors.New("Failed to find argument")
}
func (arg *PTTaskMessageArgsData) GetStringArg(name string) (string, error) {
	if val, err := arg.GetArg(name); err != nil {
		return "", err
	} else if val == nil {
		return "", nil
	} else {
		return getTypedValue[string](val)
	}
}
func (arg *PTTaskMessageArgsData) GetNumberArg(name string) (float64, error) {
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
func (arg *PTTaskMessageArgsData) GetBooleanArg(name string) (bool, error) {
	if val, err := arg.GetArg(name); err != nil {
		return false, err
	} else if val == nil {
		return false, nil
	} else {
		return getTypedValue[bool](val)
	}
}
func (arg *PTTaskMessageArgsData) GetDictionaryArg(name string) (map[string]string, error) {
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

// GetFileArg returns the file UUID that was registered with Mythic before tasking
func (arg *PTTaskMessageArgsData) GetFileArg(name string) (string, error) {
	return arg.GetStringArg(name)
}

// GetPayloadListArg returns the payload UUID that was selected from a dropdown list in the UI
func (arg *PTTaskMessageArgsData) GetPayloadListArg(name string) (string, error) {
	return arg.GetStringArg(name)
}

func (arg *PTTaskMessageArgsData) GetChooseOneArg(name string) (string, error) {
	return arg.GetStringArg(name)
}
func (arg *PTTaskMessageArgsData) GetArrayArg(name string) ([]string, error) {
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
func (arg *PTTaskMessageArgsData) GetChooseMultipleArg(name string) ([]string, error) {
	return arg.GetArrayArg(name)
}
func (arg *PTTaskMessageArgsData) GetTypedArrayArg(name string) ([][]string, error) {
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

type C2ProfileInfo struct {
	Name       string                 `json:"name" mapstructure:"name"`
	Parameters map[string]interface{} `json:"parameters" mapstructure:"parameters"`
}
type ConnectionInfo struct {
	CallbackUUID  string        `json:"callback_uuid" mapstructure:"callback_uuid"`
	AgentUUID     string        `json:"agent_uuid" mapstructure:"agent_uuid"`
	Host          string        `json:"host" mapstructure:"host"`
	C2ProfileInfo C2ProfileInfo `json:"c2_profile" mapstructure:"c2_profile"`
}

// GetConnectionInfoArg returns structured information about a new P2P connection that can be established
func (arg *PTTaskMessageArgsData) GetConnectionInfoArg(name string) (ConnectionInfo, error) {
	connectionInformation := ConnectionInfo{}
	if val, err := arg.GetArg(name); err != nil {
		return connectionInformation, err
	} else if val == nil {
		return connectionInformation, nil
	} else if err := mapstructure.Decode(val, &connectionInformation); err != nil {
		return connectionInformation, err
	} else {
		return connectionInformation, nil
	}
}

// GetLinkInfoArg returns structured information about an existing (or now dead) P2P connection
func (arg *PTTaskMessageArgsData) GetLinkInfoArg(name string) (ConnectionInfo, error) {
	connectionInformation := ConnectionInfo{}
	if val, err := arg.GetArg(name); err != nil {
		return connectionInformation, err
	} else if val == nil {
		return connectionInformation, nil
	} else if err := mapstructure.Decode(val, &connectionInformation); err != nil {
		return connectionInformation, err
	} else {
		return connectionInformation, nil
	}
}

type CredentialInfo struct {
	Realm      string `json:"realm" mapstructure:"realm"`
	Account    string `json:"account" mapstructure:"account"`
	Credential string `json:"credential" mapstructure:"credential"`
	Comment    string `json:"comment" mapstructure:"comment"`
	Type       string `json:"type" mapstructure:"type"`
}

// GetCredentialArg returns all the data about a credential from Mythic's credential store
func (arg *PTTaskMessageArgsData) GetCredentialArg(name string) (CredentialInfo, error) {
	credentialInfo := CredentialInfo{}
	if val, err := arg.GetArg(name); err != nil {
		return credentialInfo, err
	} else if val == nil {
		return credentialInfo, nil
	} else if err := mapstructure.Decode(val, &credentialInfo); err != nil {
		return credentialInfo, err
	} else {
		return credentialInfo, nil
	}
}

func getTypedValue[T any](value interface{}) (T, error) {
	switch v := value.(type) {
	case T:
		return v, nil
	default:
		var emptyResult T
		logging.LogInfo("bad type", "value", value, "type", reflect.TypeOf(value))
		return emptyResult, errors.New("bad type for value")
	}
}
func (arg *PTTaskMessageArgsData) HasArg(name string) bool {
	for _, currentArg := range arg.args {
		if currentArg.Name == name || currentArg.CLIName == name {
			return true
		}
	}
	return false
}
func (arg *PTTaskMessageArgsData) GetCommandLine() string {
	return arg.commandLine
}
func (arg *PTTaskMessageArgsData) GetRawCommandLine() string {
	return arg.rawCommandLine
}
func (arg *PTTaskMessageArgsData) GetTaskingLocation() string {
	return arg.taskingLocation
}
func (arg *PTTaskMessageArgsData) AddArg(newArg CommandParameter) error {
	// first see if newArg is in our list, if it is, update it, else add it
	if newArg.ParameterGroupInformation == nil {
		newArg.ParameterGroupInformation = []ParameterGroupInfo{
			{
				ParameterIsRequired: false,
				GroupName:           "Default",
			},
		}
		newArg.value = newArg.DefaultValue
		newArg.userSupplied = true
	}
	for i := 0; i < len(arg.args); i++ {
		if arg.args[i].Name == newArg.Name {
			// just update
			arg.args[i] = newArg
			return nil
		}
	}
	// don't have it, so add it
	arg.args = append(arg.args, newArg)
	return nil
}
func (arg *PTTaskMessageArgsData) SetArgValue(name string, value interface{}) error {
	// first see if newArg is in our list, if it is, update it, else add it
	for i := 0; i < len(arg.args); i++ {
		if arg.args[i].Name == name {
			// just update
			arg.args[i].value = value
			arg.args[i].userSupplied = true
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Failed to find arg %s", name))
}
func (arg *PTTaskMessageArgsData) GetParameterGroupName() (string, error) {
	var groupNameOptions []string
	var suppliedArgNames []string
	if len(arg.args) == 0 {
		return "Default", nil
	} else if arg.manualParameterGroupName != "" {
		return arg.manualParameterGroupName, nil
	}
	// first get a unique list of all possible groupNames for the arguments
	for _, currentArg := range arg.args {
		for _, currentGroupInfo := range currentArg.ParameterGroupInformation {
			if !helpers.StringSliceContains(groupNameOptions, currentGroupInfo.GroupName) {
				groupNameOptions = append(groupNameOptions, currentGroupInfo.GroupName)
			}
		}
	}
	// only want to look at groups based on arguments we've supplied, don't take default values into account
	for _, currentArg := range arg.args {
		if currentArg.value != nil && currentArg.userSupplied {
			suppliedArgNames = append(suppliedArgNames, currentArg.Name)
			groupNameIntersection := []string{}
			for _, currentGroupInfo := range currentArg.ParameterGroupInformation {
				if helpers.StringSliceContains(groupNameOptions, currentGroupInfo.GroupName) {
					groupNameIntersection = append(groupNameIntersection, currentGroupInfo.GroupName)
				}
			}
			groupNameOptions = groupNameIntersection
		}
	}
	if len(groupNameOptions) == 0 {
		return "", errors.New(fmt.Sprintf("Supplied arguments, %v, don't match any parameter group", suppliedArgNames))
	} else if len(groupNameOptions) == 1 {
		return groupNameOptions[0], nil
	} else {
		finalMatchingGroupNames := []string{}
		// check to see, for any possible group, if it has all of its required values, or if it's still too ambiguous
		for _, groupNameOption := range groupNameOptions {
			hasAllValues := true
			for _, currentArg := range arg.args {
				for _, currentGroupInfo := range currentArg.ParameterGroupInformation {
					if currentGroupInfo.GroupName == groupNameOption {
						if currentGroupInfo.ParameterIsRequired && !currentArg.userSupplied {
							// this parameter group has a parameter that's required that we haven't explicitly set yet
							hasAllValues = false
						}
					}
				}
			}
			if hasAllValues {
				finalMatchingGroupNames = append(finalMatchingGroupNames, groupNameOption)
			}
		}
		if len(finalMatchingGroupNames) == 0 {
			return "", errors.New(fmt.Sprintf("Supplied Arguments, %v, match more than one parameter group (%v), and all require at least one more value from the user", suppliedArgNames, groupNameOptions))
		} else if len(finalMatchingGroupNames) > 1 {
			if helpers.StringSliceContains(finalMatchingGroupNames, arg.initialParameterGroupName) {
				return arg.initialParameterGroupName, nil
			}
			return "", errors.New(fmt.Sprintf("Supplied Arguments, %v, match more than one parameter group (%v)", suppliedArgNames, finalMatchingGroupNames))
		} else {
			return finalMatchingGroupNames[0], nil
		}
	}
}
func (arg *PTTaskMessageArgsData) GetParameterGroupArguments() ([]CommandParameter, error) {
	if groupName, err := arg.GetParameterGroupName(); err != nil {
		return nil, err
	} else {
		groupArguments := []CommandParameter{}
		for _, currentArg := range arg.args {
			for _, currentGroup := range currentArg.ParameterGroupInformation {
				if currentGroup.GroupName == groupName {
					groupArguments = append(groupArguments, currentArg)
				}
			}
		}
		return groupArguments, nil
	}
}
func (arg *PTTaskMessageArgsData) RenameArg(oldName string, newName string) error {
	for i := 0; i < len(arg.args); i++ {
		if arg.args[i].Name == oldName {
			arg.args[i].Name = newName
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Failed to find argument %s", oldName))
}
func (arg *PTTaskMessageArgsData) RemoveArg(name string) error {
	for index, currentArg := range arg.args {
		if currentArg.Name == name {
			arg.args = append(arg.args[:index], arg.args[index+1:]...)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Failed to find arg %s", name))
}
func (arg *PTTaskMessageArgsData) LoadArgsFromJSONString(stringArgs string) error {
	var jsonArgs map[string]interface{}
	logging.LogDebug("Calling LoadArgsFromJSONString", "stringArgs", stringArgs)
	if err := json.Unmarshal([]byte(stringArgs), &jsonArgs); err != nil {
		logging.LogError(err, "Failed to load stringArgs into map[string]interface{}")
		return err
	} else {
		for i := 0; i < len(arg.args); i++ {
			for key, val := range jsonArgs {
				if arg.args[i].Name == key || arg.args[i].CLIName == key {
					arg.args[i].value = val
					arg.args[i].userSupplied = true
				}
			}
		}
	}
	return nil
}
func (arg *PTTaskMessageArgsData) LoadArgsFromDictionary(dictionaryArgs map[string]interface{}) error {
	//logging.LogDebug("Calling LoadArgsFromDictionary", "dictionaryArgs", dictionaryArgs, "arg.args", arg.args)
	for i := 0; i < len(arg.args); i++ {
		for key, val := range dictionaryArgs {
			logging.LogTrace("searching through dictionaryArgs for a match", "currentArg", key)
			if arg.args[i].Name == key || arg.args[i].CLIName == key {
				logging.LogTrace("Found a matching arg for LoadArgsFromDictionary", "arg", key)
				arg.args[i].value = val
				arg.args[i].userSupplied = true
			}
		}
	}
	return nil
}
func (arg *PTTaskMessageArgsData) VerifyRequiredArgsHaveValues() (bool, error) {
	if groupName, err := arg.GetParameterGroupName(); err != nil {
		return false, err
	} else {
		for _, currentArg := range arg.args {
			machedArg := false
			argRequired := false
			for _, currentGroup := range currentArg.ParameterGroupInformation {
				if currentGroup.GroupName == groupName {
					machedArg = true
					argRequired = currentGroup.ParameterIsRequired
				}
			}
			if machedArg {
				if currentArg.value == nil {
					currentArg.value = currentArg.DefaultValue
				}
				if argRequired && !currentArg.userSupplied {
					return false, errors.New(fmt.Sprintf("Required arg, %s, was not specified by the user.\nIf you did specify it, there might be an issue with the command's parsing functions", currentArg.Name))
				}
			}
		}
		return true, nil
	}

}
func (arg *PTTaskMessageArgsData) GetFinalArgs() (string, error) {
	if arg.manualArgs != nil {
		return *arg.manualArgs, nil
	} else if groupArgs, err := arg.GetParameterGroupArguments(); err != nil {
		return "", err
	} else if len(groupArgs) == 0 {
		return arg.GetCommandLine(), nil
	} else {
		// go through groupArgs and make a JSON string out of the argument values
		argMap := map[string]interface{}{}
		for _, currentArg := range groupArgs {
			currentArgValue := currentArg.GetCurrentValue()
			if currentArgValue != nil {
				argMap[currentArg.Name] = currentArg.GetCurrentValue()
			} else {
				logging.LogInfo("Not sending nil value for task", "parameter", currentArg.Name)
			}
		}
		if jsonBytes, err := json.Marshal(argMap); err != nil {
			logging.LogError(err, "Failed to convert args to JSON string", "argMap", argMap)
			return "", err
		} else {
			return string(jsonBytes), nil
		}
	}
}
func (arg *PTTaskMessageArgsData) GetUnusedArgs() string {
	if arg.manualArgs != nil {
		return "Manual args explicitly set, all args unused\n"
	}
	groupName, err := arg.GetParameterGroupName()
	if err != nil {
		return fmt.Errorf("failed to get parameter group name: %w", err).Error()
	}
	returnString := fmt.Sprintf("The following args aren't being used because they don't belong to the %s parameter group: \n",
		groupName)
	unusedGroupArguments := []CommandParameter{}
	for _, currentArg := range arg.args {
		foundGroup := false
		for _, currentGroup := range currentArg.ParameterGroupInformation {
			if currentGroup.GroupName == groupName {
				foundGroup = true
			}
		}
		if !foundGroup {
			unusedGroupArguments = append(unusedGroupArguments, currentArg)
		}
	}
	if len(unusedGroupArguments) == 0 {
		return returnString + "No arguments unused\n"
	}
	// go through groupArgs and make a JSON string out of the argument values
	argMap := map[string]interface{}{}
	for _, currentArg := range unusedGroupArguments {
		argMap[currentArg.Name] = currentArg.GetCurrentValue()
	}
	if jsonBytes, err := json.MarshalIndent(argMap, "", "  "); err != nil {
		logging.LogError(err, "Failed to convert args to JSON string", "argMap", argMap)
		return returnString + fmt.Errorf("failed to convert unused args to JSON string: %w", err).Error()
	} else {
		return returnString + string(jsonBytes)
	}
}

package agentstructs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils"
)

// Args helper functions
func GenerateArgsData(cmdParams []CommandParameter, task PTTaskMessageAllData) (PTTaskMessageArgsData, error) {
	args := PTTaskMessageArgsData{
		taskingLocation: task.Task.TaskingLocation,
		commandLine:     task.Task.Params,
		rawCommandLine:  task.Task.OriginalParams,
	}
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
	} else {
		return getTypedValue[string](val)
	}
}
func (arg *PTTaskMessageArgsData) GetNumberArg(name string) (float64, error) {
	if val, err := arg.GetArg(name); err != nil {
		return 0, err
	} else {
		return getTypedValue[float64](val)
	}
}
func (arg *PTTaskMessageArgsData) GetBooleanArg(name string) (bool, error) {
	if val, err := arg.GetArg(name); err != nil {
		return false, err
	} else {
		return getTypedValue[bool](val)
	}
}
func (arg *PTTaskMessageArgsData) GetDictionaryArg(name string) (map[string]interface{}, error) {
	if val, err := arg.GetArg(name); err != nil {
		return nil, err
	} else {
		return getTypedValue[map[string]interface{}](val)
	}
}
func getTypedValue[T any](value interface{}) (T, error) {
	switch v := value.(type) {
	case T:
		return v, nil
	default:
		var emptyResult T
		logging.LogError(nil, "Bad Type for value", "value", value)
		return emptyResult, errors.New("Bad type for value")
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
	}
	// first get a unique list of all possible groupNames for the arguments
	for _, currentArg := range arg.args {
		for _, currentGroupInfo := range currentArg.ParameterGroupInformation {
			if !utils.StringSliceContains(groupNameOptions, currentGroupInfo.GroupName) {
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
				if utils.StringSliceContains(groupNameOptions, currentGroupInfo.GroupName) {
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
			argMap[currentArg.Name] = currentArg.GetCurrentValue()
		}
		if jsonBytes, err := json.Marshal(argMap); err != nil {
			logging.LogError(err, "Failed to convert args to JSON string", "argMap", argMap)
			return "", err
		} else {
			return string(jsonBytes), nil
		}
	}
}

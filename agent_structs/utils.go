package agentstructs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"

	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/helpers"
)

// REQUIRED, Don't Modify
type allPayloadData struct {
	allCommands       []Command
	payloadDefinition PayloadType
	mutex             sync.RWMutex
	containerVersion  string
	rpcMethods        []sharedStructs.RabbitmqRPCMethod
	directMethods     []sharedStructs.RabbitmqDirectMethod
	buildFunction     func(PayloadBuildMessage) PayloadBuildResponse
}

var (
	AllPayloadData containerPayloadData
)

type containerPayloadData struct {
	PayloadMap map[string]*allPayloadData
}

func (r *containerPayloadData) GetAllPayloadTypeNames() []string {
	names := []string{}
	for key, _ := range r.PayloadMap {
		if key != "" && !helpers.StringSliceContains(names, key) {
			names = append(names, key)
		}

	}
	return names
}
func (r *containerPayloadData) Get(name string) *allPayloadData {
	if r.PayloadMap == nil {
		r.PayloadMap = make(map[string]*allPayloadData)
	}
	if existingC2Data, ok := r.PayloadMap[name]; !ok {
		newC2Data := allPayloadData{}
		r.PayloadMap[name] = &newC2Data
		return &newC2Data
	} else {
		return existingC2Data
	}
}
func (r *allPayloadData) AddCommand(cmd Command) {

	if cmd.AssociatedBrowserScript != nil {
		if cmd.AssociatedBrowserScript.ScriptPath != "" {
			if scriptContents, err := os.ReadFile(cmd.AssociatedBrowserScript.ScriptPath); err != nil {
				if err != io.EOF {
					fmt.Printf("Failed to open script file: %s, %v\n", cmd.AssociatedBrowserScript.ScriptPath, err)
					//logging.LogError(err, "Failed to open script file")
					cmd.AssociatedBrowserScript = nil
				}
			} else {
				cmd.AssociatedBrowserScript.ScriptContents = string(scriptContents)
			}
		}
	}
	for i := 0; i < len(cmd.CommandParameters); i++ {
		if cmd.CommandParameters[i].CLIName == "" {
			cmd.CommandParameters[i].CLIName = strings.ReplaceAll(cmd.CommandParameters[i].Name, " ", "-")
		}
		if len(cmd.CommandParameters[i].ParameterGroupInformation) == 0 {
			cmd.CommandParameters[i].ParameterGroupInformation = append(cmd.CommandParameters[i].ParameterGroupInformation, ParameterGroupInfo{
				GroupName:           "Default",
				ParameterIsRequired: true,
			})
		}
		for j := 0; j < len(cmd.CommandParameters[i].ParameterGroupInformation); j++ {
			if cmd.CommandParameters[i].DefaultValue == nil && !cmd.CommandParameters[i].ParameterGroupInformation[j].ParameterIsRequired {
				logging.LogWarning("default value should be set to blank value of appropriate type or a meaningful default value for this parameter instead of nil",
					"command", cmd.Name, "parameter", cmd.CommandParameters[i].Name)
			}
		}
	}
	if cmd.CommandAttributes.FilterCommandAvailabilityByAgentBuildParameters == nil {
		cmd.CommandAttributes.FilterCommandAvailabilityByAgentBuildParameters = make(map[string]string)
	}
	if cmd.CommandAttributes.SupportedOS == nil {
		cmd.CommandAttributes.SupportedOS = make([]string, 0)
	}
	for i := 0; i < len(r.allCommands); i++ {
		if r.allCommands[i].Name == cmd.Name {
			logging.LogDebug("can't add command, duplicate name detected, overwriting old one",
				"two commands with same name detected", cmd.Name)
			r.mutex.Lock()
			r.allCommands[i] = cmd
			r.mutex.Unlock()
			return
		}
	}
	r.mutex.Lock()
	r.allCommands = append(r.allCommands, cmd)
	r.mutex.Unlock()
}
func (r *allPayloadData) RemoveCommand(cmd Command) {
	for i := 0; i < len(r.allCommands); i++ {
		if r.allCommands[i].Name == cmd.Name {
			r.mutex.Lock()
			r.allCommands = append(r.allCommands[:i], r.allCommands[i+1:]...)
			r.mutex.Unlock()
			return
		}
	}
	logging.LogWarning("failed to find command for removal", "command", cmd.Name)
}
func (r *allPayloadData) AddBuildFunction(f func(PayloadBuildMessage) PayloadBuildResponse) {
	r.buildFunction = f
}
func (r *allPayloadData) AddOnNewCallbackFunction(f func(PTOnNewCallbackAllData) PTOnNewCallbackResponse) {
	r.payloadDefinition.OnNewCallback = f
}
func (r *allPayloadData) AddCheckIfCallbacksAliveFunction(f func(PTCheckIfCallbacksAliveMessage) PTCheckIfCallbacksAliveMessageResponse) {
	r.payloadDefinition.CheckIfCallbacksAliveFunction = f
}
func (r *allPayloadData) AddPayloadDefinition(payloadDef PayloadType) {
	if payloadDef.CustomRPCFunctions == nil {
		payloadDef.CustomRPCFunctions = make(map[string]func(message PTRPCOtherServiceRPCMessage) PTRPCOtherServiceRPCMessageResponse)
	}
	r.payloadDefinition = payloadDef
}
func (r *allPayloadData) AddIcon(filePath string) {
	if r.payloadDefinition.AgentIcon == nil {
		if _, err := os.Stat(filePath); err != nil {
			logging.LogError(err, "Failed to find agent icon")
			r.payloadDefinition.AgentIcon = nil
		} else if file, err := os.Open(filePath); err != nil {
			r.payloadDefinition.AgentIcon = nil
			logging.LogError(err, "Failed to open file path for agent icon")
			os.Exit(1)
		} else if agentIcon, err := io.ReadAll(file); err != nil {
			r.payloadDefinition.AgentIcon = nil
			logging.LogError(err, "Failed to read agent icon")
			os.Exit(1)
		} else {
			r.payloadDefinition.AgentIcon = &agentIcon
		}
	}
}
func (r *allPayloadData) AddDarkModeIcon(filePath string) {
	if r.payloadDefinition.DarkModeAgentIcon == nil {
		if _, err := os.Stat(filePath); err != nil {
			logging.LogError(err, "Failed to find agent icon")
			r.payloadDefinition.DarkModeAgentIcon = nil
		} else if file, err := os.Open(filePath); err != nil {
			r.payloadDefinition.DarkModeAgentIcon = nil
			logging.LogError(err, "Failed to open file path for agent icon")
			os.Exit(1)
		} else if agentIcon, err := io.ReadAll(file); err != nil {
			r.payloadDefinition.DarkModeAgentIcon = nil
			logging.LogError(err, "Failed to read agent icon")
			os.Exit(1)
		} else {
			r.payloadDefinition.DarkModeAgentIcon = &agentIcon
		}
	}
}
func (r *allPayloadData) GetPayloadDefinition() PayloadType {

	return r.payloadDefinition
}
func (r *allPayloadData) GetCommands() []Command {
	for commandIndex := range r.allCommands {
		for paramIndex := range r.allCommands[commandIndex].CommandParameters {
			for groupInfoIndex := range r.allCommands[commandIndex].CommandParameters[paramIndex].ParameterGroupInformation {
				if r.allCommands[commandIndex].CommandParameters[paramIndex].ParameterGroupInformation[groupInfoIndex].GroupName == "" {
					r.allCommands[commandIndex].CommandParameters[paramIndex].ParameterGroupInformation[groupInfoIndex].GroupName = "Default"
				}
			}
		}
	}
	return r.allCommands
}
func (r *allPayloadData) GetBuildParameters() []BuildParameter {
	return r.payloadDefinition.BuildParameters
}
func (r *allPayloadData) GetBuildFunction() func(PayloadBuildMessage) PayloadBuildResponse {
	return r.buildFunction
}
func (r *allPayloadData) GetOnNewCallbackFunction() func(PTOnNewCallbackAllData) PTOnNewCallbackResponse {
	return r.payloadDefinition.OnNewCallback
}
func (r *allPayloadData) GetCheckIfCallbacksAliveFunction() func(PTCheckIfCallbacksAliveMessage) PTCheckIfCallbacksAliveMessageResponse {
	return r.payloadDefinition.CheckIfCallbacksAliveFunction
}
func (r *allPayloadData) AddContainerVersion(ver string) {
	r.containerVersion = ver
}
func (r *allPayloadData) GetPayloadName() string {
	return r.payloadDefinition.Name
}
func (r *allPayloadData) GetContainerVersion() string {
	return r.containerVersion
}
func (r *allPayloadData) AddRPCMethod(m sharedStructs.RabbitmqRPCMethod) {
	r.mutex.Lock()
	r.rpcMethods = append(r.rpcMethods, m)
	r.mutex.Unlock()
}
func (r *allPayloadData) GetRPCMethods() []sharedStructs.RabbitmqRPCMethod {
	return r.rpcMethods
}
func (r *allPayloadData) AddDirectMethod(m sharedStructs.RabbitmqDirectMethod) {
	r.mutex.Lock()
	r.directMethods = append(r.directMethods, m)
	r.mutex.Unlock()
}
func (r *allPayloadData) GetDirectMethods() []sharedStructs.RabbitmqDirectMethod {
	return r.directMethods
}
func (r *allPayloadData) GetRoutingKey(baseRoutingKey string) string {
	return fmt.Sprintf("%s_%s", r.GetPayloadName(), baseRoutingKey)
}
func RunShellCommand(arguments string, cwd string) (stdout []byte, stderr []byte, err error) {
	return RunCommand("/bin/bash", arguments, cwd)
}
func RunCommand(command string, arguments string, cwd string) (stdout []byte, stderr []byte, err error) {
	cmd := exec.Command(command)
	cmd.Stdin = strings.NewReader(arguments)
	cmd.Dir = cwd
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	errOut := cmd.Run()
	return stdOut.Bytes(), stdErr.Bytes(), errOut
}

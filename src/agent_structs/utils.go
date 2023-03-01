package agentstructs

import (
	"bytes"
	"context"
	"fmt"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type RabbitmqRPCMethod struct {
	RabbitmqRoutingKey         string
	RabbitmqProcessingFunction func([]byte) interface{}
}
type RabbitmqDirectMethod struct {
	RabbitmqRoutingKey         string
	RabbitmqProcessingFunction func([]byte)
}

// REQUIRED, Don't Modify
type allPayloadData struct {
	allCommands       []Command
	payloadDefinition PayloadType
	mutex             sync.RWMutex
	containerVersion  string
	rpcMethods        []RabbitmqRPCMethod
	directMethods     []RabbitmqDirectMethod
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
		if key != "" && !utils.StringSliceContains(names, key) {
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
	}
	if cmd.CommandAttributes.FilterCommandAvailabilityByAgentBuildParameters == nil {
		cmd.CommandAttributes.FilterCommandAvailabilityByAgentBuildParameters = make(map[string]string)
	}
	if cmd.CommandAttributes.SupportedOS == nil {
		cmd.CommandAttributes.SupportedOS = make([]string, 0)
	}
	r.mutex.Lock()
	r.allCommands = append(r.allCommands, cmd)
	r.mutex.Unlock()
}
func (r *allPayloadData) AddBuildFunction(f func(PayloadBuildMessage) PayloadBuildResponse) {
	r.buildFunction = f
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
		} else if agentIcon, err := ioutil.ReadAll(file); err != nil {
			r.payloadDefinition.AgentIcon = nil
			logging.LogError(err, "Failed to read agent icon")
			os.Exit(1)
		} else {
			r.payloadDefinition.AgentIcon = &agentIcon
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
func (r *allPayloadData) GetBuildFunction() func(PayloadBuildMessage) PayloadBuildResponse {
	return r.buildFunction
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
func (r *allPayloadData) AddRPCMethod(m RabbitmqRPCMethod) {
	r.mutex.Lock()
	r.rpcMethods = append(r.rpcMethods, m)
	r.mutex.Unlock()
}
func (r *allPayloadData) GetRPCMethods() []RabbitmqRPCMethod {
	return r.rpcMethods
}
func (r *allPayloadData) AddDirectMethod(m RabbitmqDirectMethod) {
	r.mutex.Lock()
	r.directMethods = append(r.directMethods, m)
	r.mutex.Unlock()
}
func (r *allPayloadData) GetDirectMethods() []RabbitmqDirectMethod {
	return r.directMethods
}
func (r *allPayloadData) GetRoutingKey(baseRoutingKey string) string {
	return fmt.Sprintf("%s_%s", r.GetPayloadName(), baseRoutingKey)
}

func RunCommandWithTimeout(command string, args []string, cwd string, timeoutSeconds int) (stdout []byte, stderr []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = cwd
	var stdOutBytes bytes.Buffer
	var stdErrBytes bytes.Buffer
	cmd.Stdout = &stdOutBytes
	cmd.Stderr = &stdErrBytes
	if err := cmd.Start(); err != nil {
		logging.LogError(err, "Failed to run command")
		return nil, nil, err
	} else if err := cmd.Wait(); err != nil {
		logging.LogError(err, "Command failed to complete successfully")
		return nil, nil, err
	} else {
		return stdOutBytes.Bytes(), stdErrBytes.Bytes(), nil
	}
}

package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         PT_TASK_CREATE_TASKING,
		RabbitmqProcessingFunction: processPtTaskCreateMessages,
	})
}

type ptTaskCreateTaskingMessageResponseWrapper struct {
	agentstructs.PTTaskCreateTaskingMessageResponse
	Params string `json:"params"`
}

func processPtTaskCreateMessages(msg []byte) {
	incomingMessage := agentstructs.PTTaskMessageAllData{}
	response := ptTaskCreateTaskingMessageResponseWrapper{}
	if err := json.Unmarshal(msg, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Success = false
		response.Error = "Failed to unmarshal JSON message into structs"
		sendTaskCreateResponse(response)
		return
	} else {
		response.Success = false
		response.TaskID = incomingMessage.Task.ID
		for _, command := range agentstructs.AllPayloadData.Get(incomingMessage.CommandPayloadType).GetCommands() {
			if command.Name == incomingMessage.Task.CommandName {
				if err := prepTaskArgs(command, &incomingMessage); err != nil {
					response.Success = false
					response.Error = err.Error()
					sendTaskCreateResponse(response)
					return
				}
				agentResponse := command.TaskFunctionCreateTasking(&incomingMessage)
				finalResponse := ptTaskCreateTaskingMessageResponseWrapper{
					agentResponse,
					"",
				}
				if finalResponse.Success {
					if incomingMessage.Task.IsInteractiveTask {
						finalResponse.Params = incomingMessage.Args.GetFinalInteractiveTaskingArgs()
					} else if requiredArgsHaveValues, err := incomingMessage.Args.VerifyRequiredArgsHaveValues(); err != nil {
						logging.LogError(err, "Failed to verify if all required args have values")
						finalResponse.Success = false
						finalResponse.Error = fmt.Sprintf("Failed to verify if all required args have values:\n%s", err.Error())
					} else if !requiredArgsHaveValues {
						finalResponse.Success = false
						finalResponse.Error = fmt.Sprintf("Some required args are missing values")
					} else if params, err := incomingMessage.Args.GetFinalArgs(); err != nil {
						logging.LogError(err, "Failed to get final arguments", "args", incomingMessage.Args)
						finalResponse.Success = false
						finalResponse.Error = fmt.Sprintf("Failed to generate final arguments:\n%s", err.Error())
					} else {
						finalResponse.Params = params
						unusedParams := incomingMessage.Args.GetUnusedArgs()
						if finalResponse.Stdout != nil {
							*finalResponse.Stdout = (*finalResponse.Stdout) + "\n" + unusedParams
						} else {
							finalResponse.Stdout = &unusedParams
						}
						if finalResponse.ParameterGroupName == "" {
							if newGroupName, err := incomingMessage.Args.GetParameterGroupName(); err != nil {
								logging.LogError(err, "Failed to get new parameter group name for task")
								finalResponse.Success = false
								finalResponse.Error = err.Error()
							} else {
								finalResponse.ParameterGroupName = newGroupName
							}
						}
					}
				}
				sendTaskCreateResponse(finalResponse)
				return
			}
		}
		// if we get here then we never found the command
		response.Error = fmt.Sprintf("Failed to find command: %s", incomingMessage.Task.CommandName)
		sendTaskCreateResponse(response)
		return
	}
}

func sendTaskCreateResponse(response ptTaskCreateTaskingMessageResponseWrapper) {
	for {
		err := RabbitMQConnection.SendStructMessage(
			MYTHIC_EXCHANGE,
			PT_TASK_CREATE_TASKING_RESPONSE,
			"",
			response,
			false,
		)
		if err != nil {
			logging.LogError(err, "Failed to send payload response back to Mythic")
			continue
		}
		return
	}
}

package rabbitmq

import (
	"encoding/json"
	"fmt"

	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
)

func init() {
	agentstructs.AllPayloadData.Get("").AddDirectMethod(agentstructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         PT_TASK_CREATE_TASKING,
		RabbitmqProcessingFunction: processPtTaskCreateMessages,
	})
}

func processPtTaskCreateMessages(msg []byte) {
	incomingMessage := agentstructs.PTTaskMessageAllData{}
	response := agentstructs.PTTaskCreateTaskingMessageResponse{}
	if err := json.Unmarshal(msg, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		response.Success = false
		response.Error = "Failed to unmarshal JSON message into structs"
		sendTaskCreateResponse(response)
		return
	} else {
		response.Success = false
		response.TaskID = incomingMessage.Task.ID
		for _, command := range agentstructs.AllPayloadData.Get(incomingMessage.PayloadType).GetCommands() {
			if command.Name == incomingMessage.Task.CommandName {
				if err := prepTaskArgs(command, &incomingMessage); err != nil {
					response.Success = false
					response.Error = err.Error()
					sendTaskCreateResponse(response)
					return
				}
				response = command.TaskFunctionCreateTasking(incomingMessage)
				if response.Success {
					if requiredArgsHaveValues, err := incomingMessage.Args.VerifyRequiredArgsHaveValues(); err != nil {
						logging.LogError(err, "Failed to verify if all required args have values")
						response.Success = false
						response.Error = fmt.Sprintf("Failed to verify if all required args have values:\n%s", err.Error())
					} else if !requiredArgsHaveValues {
						response.Success = false
						response.Error = fmt.Sprintf("Some required args are missing values")
					} else if params, err := incomingMessage.Args.GetFinalArgs(); err != nil {
						logging.LogError(err, "Failed to get final arguments", "args", incomingMessage.Args)
						response.Success = false
						response.Error = fmt.Sprintf("Failed to generate final arguments:\n%s", err.Error())
					} else {
						response.Params = params
						if response.ParameterGroupName == "" {
							if newGroupName, err := incomingMessage.Args.GetParameterGroupName(); err != nil {
								logging.LogError(err, "Failed to get new parameter group name for task")
								response.Success = false
								response.Error = err.Error()
							} else {
								response.ParameterGroupName = newGroupName
							}
						}
					}
				}
				sendTaskCreateResponse(response)
				return
			}
		}
		// if we get here then we never found the command
		response.Error = fmt.Sprintf("Failed to find command: %s", incomingMessage.Task.CommandName)
		sendTaskCreateResponse(response)
		return
	}
}

func sendTaskCreateResponse(response agentstructs.PTTaskCreateTaskingMessageResponse) {
	if err := RabbitMQConnection.SendStructMessage(
		MYTHIC_EXCHANGE,
		PT_TASK_CREATE_TASKING_RESPONSE,
		"",
		response,
	); err != nil {
		logging.LogError(err, "Failed to send payload response back to Mythic")
	}
}

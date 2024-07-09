package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"time"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_GET_SERVER_DEBUG_OUTPUT,
		RabbitmqProcessingFunction: processC2RPCGetDebugOutput,
	})
}
func processC2RPCGetDebugOutput(msg []byte) interface{} {
	input := c2structs.C2GetDebugOutputMessage{}
	responseMsg := c2structs.C2GetDebugOutputMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually start the c2 profile if needed
		return C2RPCGetDebugOutput(input)
	}
	return responseMsg
}

func C2RPCGetDebugOutput(input c2structs.C2GetDebugOutputMessage) c2structs.C2GetDebugOutputMessageResponse {
	responseMsg := c2structs.C2GetDebugOutputMessageResponse{
		Success: false,
		Error:   "Not implemented, not getting debug output",
	}
	if c2structs.AllC2Data.Get(input.Name).RunningServerProcess != nil {
		output := ""
		finishedReadingOutput := make(chan bool, 1)
		tellGoroutineToFinish := make(chan bool, 1)
		go func() {
			<-time.After(3 * time.Second)
			tellGoroutineToFinish <- true
		}()
		go func() {
			for {
				select {
				case <-tellGoroutineToFinish:
					return
				case newOutput, ok := <-c2structs.AllC2Data.Get(input.Name).OutputChannel:
					if !ok {
						finishedReadingOutput <- true
						return
					} else {
						output += newOutput + "\n"
					}
				}
			}
		}()
		<-finishedReadingOutput

		if c2structs.AllC2Data.Get(input.Name).RunningServerProcess != nil && c2structs.AllC2Data.Get(input.Name).RunningServerProcess.ProcessState.ExitCode() == -1 {
			// we're still running
			responseMsg.Message = output
			if responseMsg.Message == "" {
				responseMsg.Message = "No Server Output\n"
			}
			responseMsg.Success = true
		} else {
			err := c2structs.AllC2Data.Get(input.Name).RunningServerProcess.Wait()
			if err != nil {
				responseMsg.Message = "Process died with error: " + err.Error()
			} else {
				responseMsg.Message = "Process exited without error"
			}

		}
	} else {
		responseMsg.Message = "Server not running\n"
	}
	return responseMsg
}

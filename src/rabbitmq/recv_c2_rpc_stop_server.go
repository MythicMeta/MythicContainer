package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"time"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(c2structs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_STOP_SERVER_ROUTING_KEY,
		RabbitmqProcessingFunction: processC2RPCStopServer,
	})
}

func processC2RPCStopServer(msg []byte) interface{} {
	input := c2structs.C2RPCStopServerMessage{}
	responseMsg := c2structs.C2RPCStopServerMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return C2RPCStopServer(input)
	}
	return responseMsg
}

func C2RPCStopServer(input c2structs.C2RPCStopServerMessage) c2structs.C2RPCStopServerMessageResponse {
	responseMsg := c2structs.C2RPCStopServerMessageResponse{
		Success: false,
		Error:   "Not implemented, not stopping",
	}
	if c2structs.AllC2Data.Get(input.Name).RunningServerProcess == nil {
		responseMsg.Error = "Server not running"
		responseMsg.InternalServerRunning = false
		return responseMsg
	} else if err := c2structs.AllC2Data.Get(input.Name).RunningServerProcess.Process.Kill(); err != nil {
		responseMsg.Error = err.Error()
		responseMsg.InternalServerRunning = false
		c2structs.AllC2Data.Get(input.Name).RunningServerProcess = nil
		return responseMsg
	} else {
		c2structs.AllC2Data.Get(input.Name).RunningServerProcess.Process.Wait()
		responseMsg.Error = ""
		output := ""
		finishedReadingOutput := make(chan bool, 1)
		tellGoroutineToFinish := make(chan bool, 1)
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
				case <-time.After(3 * time.Second):
					finishedReadingOutput <- true
					return
				}

			}
			finishedReadingOutput <- true
		}()
		select {
		case <-finishedReadingOutput:
			tellGoroutineToFinish <- true
		case <-time.After(3 * time.Second):
			tellGoroutineToFinish <- true
		}
		responseMsg.Message = output
		responseMsg.InternalServerRunning = false
		responseMsg.Success = true
		c2structs.AllC2Data.Get(input.Name).RunningServerProcess = nil
		return responseMsg
	}
}

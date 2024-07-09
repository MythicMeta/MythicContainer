package rabbitmq

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_START_SERVER_ROUTING_KEY,
		RabbitmqProcessingFunction: processC2RPCStartServer,
	})
}

func processC2RPCStartServer(msg []byte) interface{} {
	input := c2structs.C2RPCStartServerMessage{}
	responseMsg := c2structs.C2RPCStartServerMessageResponse{
		Success: false,
		Error:   "Not implemented, not starting",
	}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return C2RPCStartServer(input)
	}
	return responseMsg
}

func C2RPCStartServer(input c2structs.C2RPCStartServerMessage) c2structs.C2RPCStartServerMessageResponse {
	responseMsg := c2structs.C2RPCStartServerMessageResponse{
		Success: false,
		Error:   "",
	}
	if c2structs.AllC2Data.Get(input.Name).RunningServerProcess != nil {
		responseMsg.Error = fmt.Sprintf("Server already running with pid %d", c2structs.AllC2Data.Get(input.Name).RunningServerProcess.Process.Pid)
		responseMsg.Success = true
		responseMsg.InternalServerRunning = true
		return responseMsg
	}
	if serverFilePath, err := filepath.Abs(c2structs.AllC2Data.Get(input.Name).GetServerPath()); err != nil {
		logging.LogError(err, "Failed to get absolute path to server binary")
		responseMsg.Error = fmt.Sprintf("Failed to get absolute path to binary: %s\n", c2structs.AllC2Data.Get(input.Name).GetServerPath())
		return responseMsg
	} else if _, err := os.Stat(serverFilePath); err != nil {
		logging.LogError(err, "Failed to stat server file")
		responseMsg.Error = fmt.Sprintf("Failed to stat the server binary: %s\n", serverFilePath)
		return responseMsg
	} else {
		cmd := exec.Command(serverFilePath)
		cmd.Env = os.Environ()
		cmd.Dir = filepath.Dir(serverFilePath)
		stdOutPipe, err := cmd.StdoutPipe()
		if err != nil {
			logging.LogError(err, "Failed to get stdout pipe")
			responseMsg.Error = "Failed to get stdout pipe"
			return responseMsg
		}
		stdErrPipe, err := cmd.StderrPipe()
		if err != nil {
			logging.LogError(err, "Failed to get stderr pipe")
			responseMsg.Error = "Failed to get stderr pipe"
			return responseMsg
		}
		err = cmd.Start()
		if err != nil {
			logging.LogError(err, "Failed to start server sub process")
			responseMsg.Error = err.Error()
			c2structs.AllC2Data.Get(input.Name).RunningServerProcess = nil
			return responseMsg
		}
		c2structs.AllC2Data.Get(input.Name).RunningServerProcess = cmd
		go readStdOutToChan(input.Name, bufio.NewScanner(stdOutPipe))
		go readStdErrToChan(input.Name, bufio.NewScanner(stdErrPipe))

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
					finishedReadingOutput <- true
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
		result := make(chan error, 1)
		go func() {
			result <- c2structs.AllC2Data.Get(input.Name).RunningServerProcess.Wait()
		}()
		select {
		case <-time.After(1 * time.Second):
			err = nil
		case err = <-result:
		}
		if err != nil {
			responseMsg.Error = output
			c2structs.AllC2Data.Get(input.Name).RunningServerProcess = nil
			return responseMsg
		}
		responseMsg.Message = output
		responseMsg.Error = ""
		responseMsg.Success = true
		responseMsg.InternalServerRunning = true
		return responseMsg
	}
}

func readStdOutToChan(name string, stdOut *bufio.Scanner) {
	for stdOut.Scan() {
		output := stdOut.Text()
		logging.LogDebug(output)
		select {
		case c2structs.AllC2Data.Get(name).OutputChannel <- output:
		default:
		}
	}
	logging.LogDebug("readStdOutToChan exited")
}
func readStdErrToChan(name string, stdErr *bufio.Scanner) {
	for stdErr.Scan() {
		output := stdErr.Text()
		logging.LogDebug(output)
		select {
		case c2structs.AllC2Data.Get(name).OutputChannel <- output:
		default:
		}
	}
	logging.LogDebug("readStdErrToChan exited")
}

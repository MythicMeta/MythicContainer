package rabbitmq

import (
	"encoding/json"
	"github.com/MythicMeta/MythicContainer/c2_structs"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	c2structs.AllC2Data.Get("").AddRPCMethod(sharedStructs.RabbitmqRPCMethod{
		RabbitmqRoutingKey:         C2_RPC_REDIRECTOR_RULES_ROUTING_KEY,
		RabbitmqProcessingFunction: processC2RPCGetRedirectorRules,
	})
}

func processC2RPCGetRedirectorRules(msg []byte) interface{} {
	input := c2structs.C2GetRedirectorRuleMessage{}
	responseMsg := c2structs.C2GetRedirectorRuleMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		return C2RPCGetRedirectorRules(input)
	}
	return responseMsg
}

func C2RPCGetRedirectorRules(input c2structs.C2GetRedirectorRuleMessage) c2structs.C2GetRedirectorRuleMessageResponse {
	responseMsg := c2structs.C2GetRedirectorRuleMessageResponse{
		Success: false,
		Error:   "Not implemented, not getting redirector rules",
	}
	c2Mutex.Lock()
	if c2structs.AllC2Data.Get(input.Name).GetC2Definition().GetRedirectorRulesFunction != nil {
		responseMsg = c2structs.AllC2Data.Get(input.Name).GetC2Definition().GetRedirectorRulesFunction(input)
	}
	c2Mutex.Unlock()
	if responseMsg.RestartInternalServer {
		go restartC2Server(input.Name)
	}
	return responseMsg
}

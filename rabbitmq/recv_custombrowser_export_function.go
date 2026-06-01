package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MythicMeta/MythicContainer/custombrowserstructs"
	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"

	"github.com/MythicMeta/MythicContainer/logging"
)

// Register this RPC method with rabbitmq so it can be called
func init() {
	custombrowserstructs.AllCustomBrowserData.Get("").AddDirectMethod(sharedStructs.RabbitmqDirectMethod{
		RabbitmqRoutingKey:         CUSTOMBROWSER_EXPORT_FUNCTION,
		RabbitmqProcessingFunction: processCustomBrowserExportFunction,
	})
}

// All rabbitmq methods must take byte inputs and return an interface.
// However, we can cast these to the input and return types defined in this file
func processCustomBrowserExportFunction(ctx context.Context, msg []byte) {
	input := custombrowserstructs.ExportFunctionMessage{}
	responseMsg := custombrowserstructs.ExportFunctionMessageResponse{}
	if err := json.Unmarshal(msg, &input); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Success = false
		responseMsg.Error = "Failed to unmarshal JSON message into structs"
	} else {
		// actually do config checks on configCheck
		responseMsg = CustomBrowserExportFunction(ctx, input)
	}
	for {
		err := RabbitMQConnection.SendStructMessageWithContext(ctx,
			MYTHIC_EXCHANGE,
			CUSTOMBROWSER_EXPORT_FUNCTION_RESPONSE,
			"",
			responseMsg,
			false,
		)
		if err != nil {
			logging.LogError(err, "Failed to send custom browser export response back to Mythic")
			time.Sleep(5 * time.Second)
			continue
		}
		logging.LogDebug("Finished processing custom browser export message")
		return
	}
}

func CustomBrowserExportFunction(ctx context.Context, input custombrowserstructs.ExportFunctionMessage) custombrowserstructs.ExportFunctionMessageResponse {
	responseMsg := custombrowserstructs.ExportFunctionMessageResponse{
		Success:     false,
		Error:       "No Export Function exists",
		OperationID: input.OperatorID,
		TreeType:    input.TreeType,
	}
	if custombrowserstructs.AllCustomBrowserData.Get(input.ContainerName).GetCustomBrowserDefinition().ExportFunction != nil {
		responseMsg = custombrowserstructs.AllCustomBrowserData.Get(input.ContainerName).GetCustomBrowserDefinition().ExportFunction(ctx, input)
		responseMsg.OperationID = input.OperationID
		responseMsg.TreeType = input.TreeType
	}
	return responseMsg
}

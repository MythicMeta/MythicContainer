package mythicrpc

import "github.com/MythicMeta/MythicContainer/utils/mythicutils"

type MythicRPCPayloadGetContentMessage struct {
	PayloadUUID string `json:"uuid"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCPayloadGetContentMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Content []byte `json:"content"`
}

func SendMythicRPCPayloadGetContent(input MythicRPCPayloadGetContentMessage) (*MythicRPCPayloadGetContentMessageResponse, error) {
	response := MythicRPCPayloadGetContentMessageResponse{}
	if contents, err := mythicutils.GetFileFromMythic(input.PayloadUUID); err != nil {
		response.Error = err.Error()
		response.Success = false
	} else {
		response.Success = true
		response.Content = *contents
	}
	return &response, nil
	/*
		if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
			rabbitmq.MYTHIC_EXCHANGE,
			rabbitmq.MYTHIC_RPC_PAYLOAD_GET_PAYLOAD_CONTENT,
			input,
		); err != nil {
			logging.LogError(err, "Failed to send RPC message")
			return nil, err
		} else if err := json.Unmarshal(responseBytes, &response); err != nil {
			logging.LogError(err, "Failed to parse response back to struct", "response", response)
			return nil, err
		} else {
			return &response, nil
		}

	*/
}

package mythicrpc

type MythicRPCFileGetContentMessage struct {
	AgentFileID string `json:"file_id"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCFileGetContentMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Content []byte `json:"content"`
}

func SendMythicRPCFileGetContent(input MythicRPCFileGetContentMessage) (*MythicRPCFileGetContentMessageResponse, error) {
	response := MythicRPCFileGetContentMessageResponse{}
	if content, err := getFileFromMythic(input.AgentFileID); err != nil {
		response.Success = false
		response.Error = err.Error()
		return &response, nil
	} else {
		response.Success = true
		response.Content = *content
		return &response, nil
	}
	/*
		if responseBytes, err := rabbitmq.RabbitMQConnection.SendRPCStructMessage(
			rabbitmq.MYTHIC_EXCHANGE,
			rabbitmq.MYTHIC_RPC_FILE_GET_CONTENT,
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

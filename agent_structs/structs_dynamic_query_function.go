package agentstructs

type PTRPCDynamicQueryFunctionMessage struct {
	Command       string `json:"command" binding:"required"`
	ParameterName string `json:"parameter_name" binding:"required"`
	PayloadType   string `json:"payload_type" binding:"required"`
	Callback      int    `json:"callback" binding:"required"`
}

type PTRPCDynamicQueryFunctionMessageResponse struct {
	Success bool     `json:"success"`
	Error   string   `json:"error"`
	Choices []string `json:"choices"`
}

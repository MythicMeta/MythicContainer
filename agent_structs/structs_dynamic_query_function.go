package agentstructs

type PTRPCDynamicQueryFunctionMessage struct {
	// Command - the command name for the query function called
	Command string `json:"command" binding:"required"`
	// ParameterName - the specific parameter for the query function called
	ParameterName string `json:"parameter_name" binding:"required"`
	// PayloadType - the name of the payload type for the query function called
	PayloadType string `json:"payload_type" binding:"required"`
	// Callback - the ID of the callback where this query function is called
	Callback int `json:"callback" binding:"required"`
}

type PTRPCDynamicQueryFunctionMessageResponse struct {
	// Success - indicating if the query function succeeded or not
	Success bool `json:"success"`
	// Error - if there was an error, return that message here for the user
	Error string `json:"error"`
	// Choices - the resulting choices for the user based on the dynamic query function
	Choices []string `json:"choices"`
}

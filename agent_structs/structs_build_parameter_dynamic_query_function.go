package agentstructs

type PTRPCDynamicQueryBuildParameterFunctionMessage struct {
	// ParameterName - the specific parameter for the query function called
	ParameterName string `json:"parameter_name" binding:"required"`
	// PayloadType - the name of the payload type of the callback for the query function called
	PayloadType string `json:"payload_type" binding:"required"`
	// SelectedOS - the string OS selected during payload creation
	SelectedOS string `json:"selected_os"`
	// Secrets - User supplied secrets
	Secrets map[string]interface{} `json:"secrets"`
}
type PTRPCDynamicQueryBuildParameterFunctionMessageResponse struct {
	// Success - indicating if the query function succeeded or not
	Success bool `json:"success"`
	// Error - if there was an error, return that message here for the user
	Error string `json:"error"`
	// Choices - the resulting choices for the user based on the dynamic query function
	Choices []string `json:"choices"`
}

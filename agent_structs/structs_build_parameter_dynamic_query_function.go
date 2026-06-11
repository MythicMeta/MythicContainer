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
	// OtherParameters - other user supplied parameters
	OtherParameters map[string]interface{} `json:"other_parameters"`
}
type ComplexChoice struct {
	DisplayValue string `json:"display_value"`
	Value        string `json:"value"`
}
type PTRPCDynamicQueryBuildParameterFunctionMessageResponse struct {
	// Success - indicating if the query function succeeded or not
	Success bool `json:"success"`
	// Error - if there was an error, return that message here for the user
	Error string `json:"error"`
	// Choices - the resulting choices for the user based on the dynamic query function
	Choices []string `json:"choices"`
	// ComplexChoices - the ability to specify a value and display value for more complex usability
	// this is the required format for options provided for a parameter of type JSONString
	ComplexChoices []ComplexChoice `json:"complex_choices"`
}

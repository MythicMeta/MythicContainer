package agentstructs

type PTRPCTypedArrayParseFunctionMessage struct {
	// Command - the command name for the query function called
	Command string `json:"command" binding:"required"`
	// ParameterName - the specific parameter for the query function called
	ParameterName string `json:"parameter_name" binding:"required"`
	// PayloadType - the name of the payload type for the callback where query function called
	PayloadType string `json:"payload_type" binding:"required"`
	// CommandPayloadType - the name of the payload type for the command issued
	CommandPayloadType string `json:"command_payload_type"`
	// Callback - the ID of the callback where this query function is called
	Callback int `json:"callback" binding:"required"`
	// InputArray - the structured input array that the user provided
	InputArray []string `json:"input_array"`
}

type PTRPCTypedArrayParseMessageResponse struct {
	// Success - indicating if the query function succeeded or not
	Success bool `json:"success"`
	// Error - if there was an error, return that message here for the user
	Error string `json:"error"`
	// TypedArray - the resulting typed array based on the formatted normal array
	TypedArray [][]string `json:"typed_array"`
}

package agentstructs

type PTRPCCommandHelpFunctionMessage struct {
	// CommandNames - The list of commands to get help for
	CommandNames []string `json:"command_names" binding:"required"`
	// PayloadType - The payload type where to task for help
	PayloadType string `json:"payload_type" binding:"required"`
}

type PTRPCCommandHelpFunctionMessageResponse struct {
	// Success - indicating if the query function succeeded or not
	Success bool `json:"success"`
	// Error - if there was an error, return that message here for the user
	Error string `json:"error"`
	// Output - The help output to show the user
	Output string `json:"output"`
}

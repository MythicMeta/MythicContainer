package agentstructs

type PTRPCDynamicQueryFunctionMessage struct {
	// Command - the command name for the query function called
	Command string `json:"command" binding:"required"`
	// ParameterName - the specific parameter for the query function called
	ParameterName string `json:"parameter_name" binding:"required"`
	// PayloadType - the name of the payload type of the callback for the query function called
	PayloadType string `json:"payload_type" binding:"required"`
	// CommandPayloadType - the name of the payload type associated with this command
	CommandPayloadType string `json:"command_payload_type"`
	// Callback - the ID of the callback where this query function is called
	Callback int `json:"callback" binding:"required"`
	// PayloadOS - the string OS selected during payload creation
	PayloadOS string `json:"payload_os"`
	// PayloadUUID - the UUID of the backing payload that can be used to fetch more information about the payload
	PayloadUUID string `json:"payload_uuid"`
	// CallbackDisplayID - the number seen on the active callbacks page for the callback in question
	CallbackDisplayID int `json:"callback_display_id"`
	// AgentCallbackID - the UUID of the callback known by the agent
	AgentCallbackID string `json:"agent_callback_id"`
	// Secrets - User supplied secrets
	Secrets map[string]interface{} `json:"secrets"`
}

type PTRPCDynamicQueryFunctionMessageResponse struct {
	// Success - indicating if the query function succeeded or not
	Success bool `json:"success"`
	// Error - if there was an error, return that message here for the user
	Error string `json:"error"`
	// Choices - the resulting choices for the user based on the dynamic query function
	Choices []string `json:"choices"`
}

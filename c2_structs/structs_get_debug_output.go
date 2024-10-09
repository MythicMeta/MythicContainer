package c2structs

// C2_GET_DEBUG_OUTPUT STRUCTS

type C2GetDebugOutputMessage struct {
	Name string `json:"c2_profile_name"`
}

type C2GetDebugOutputMessageResponse struct {
	Success               bool   `json:"success"`
	Error                 string `json:"error"`
	Message               string `json:"message"`
	InternalServerRunning bool   `json:"server_running"`
	RestartInternalServer bool   `json:"restart_internal_server"`
}

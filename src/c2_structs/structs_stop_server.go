package c2structs

// C2_STOP_SERVER STRUCTS

type C2RPCStopServerMessage struct {
	Name string `json:"c2_profile_name"`
}

type C2RPCStopServerMessageResponse struct {
	Success               bool   `json:"success"`
	Error                 string `json:"error"`
	Message               string `json:"message"`
	InternalServerRunning bool   `json:"server_running"`
}

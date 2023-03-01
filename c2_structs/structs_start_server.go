package c2structs

// C2_START_SERVER STRUCTS

type C2RPCStartServerMessage struct {
	Name string `json:"c2_profile_name"`
}

type C2RPCStartServerMessageResponse struct {
	Success               bool   `json:"success"`
	Error                 string `json:"error"`
	Message               string `json:"message"`
	InternalServerRunning bool   `json:"server_running"`
}

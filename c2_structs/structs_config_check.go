package c2structs

// C2_CONFIG_CHECK STRUCTS

type C2ConfigCheckMessage struct {
	C2Parameters
}

type C2ConfigCheckMessageResponse struct {
	Success               bool   `json:"success"`
	Error                 string `json:"error"`
	Message               string `json:"message"`
	RestartInternalServer bool   `json:"restart_internal_server"`
}

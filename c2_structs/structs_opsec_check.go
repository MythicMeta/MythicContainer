package c2structs

// C2_OPSEC_CHECKS STRUCTS

type C2OPSECMessage struct {
	C2Parameters
}

type C2OPSECMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

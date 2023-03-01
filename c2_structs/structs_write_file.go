package c2structs

// C2_WRITE_FILE STRUCTS

type C2RPCWriteFileMessage struct {
	Name     string `json:"c2_profile_name"`
	Filename string `json:"filename"`
	Contents []byte `json:"contents"`
}

type C2RPCWriteFileMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

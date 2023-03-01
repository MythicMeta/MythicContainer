package c2structs

// C2_READ_FILE STRUCTS

type C2RPCGetFileMessage struct {
	Name     string `json:"c2_profile_name"`
	Filename string `json:"filename"`
}

type C2RPCGetFileMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message []byte `json:"message"`
}

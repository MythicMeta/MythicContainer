package c2structs

type C2RPCRemoveFileMessage struct {
	Name     string `json:"c2_profile_name"`
	Filename string `json:"filename"`
}

type C2RPCRemoveFileMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

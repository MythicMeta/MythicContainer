package c2structs

// C2_READ_FILE STRUCTS

type C2RPCListFileMessage struct {
	Name string `json:"c2_profile_name"`
}

type C2RPCListFileMessageResponse struct {
	Success bool     `json:"success"`
	Error   string   `json:"error"`
	Files   []string `json:"files"`
}

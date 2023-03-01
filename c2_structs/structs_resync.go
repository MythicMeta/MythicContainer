package c2structs

// C2_RESYNC STRUCTS

type C2RPCReSyncMessage struct {
	Name string `json:"c2_profile_name"`
}

type C2RPCReSyncMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

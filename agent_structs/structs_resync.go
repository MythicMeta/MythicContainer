package agentstructs

// PT_RESYNC STRUCTS

type PTRPCReSyncMessage struct {
	Name string `json:"payload_type"`
}

type PTRPCReSyncMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

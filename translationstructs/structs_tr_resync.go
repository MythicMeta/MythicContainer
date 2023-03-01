package translationstructs

type TRRPCReSyncMessage struct {
	Name string `json:"translation_name"`
}

type TRRPCReSyncMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

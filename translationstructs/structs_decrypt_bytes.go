package translationstructs

// TRANSLATION_CONTAINER_DECRYPT_BYTES STRUCTS

type TrDecryptBytesMessage struct {
	TranslationContainerName string `json:"translation_container_name"`
	EncryptionKey            []byte `json:"enc_key"`
	CryptoType               string `json:"crypto_type"`
	Message                  []byte `json:"message"`
	AgentCallbackUUID        string `json:"callback_uuid"`
}

type TrDecryptBytesMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message []byte `json:"message"`
}

package translationstructs

// TRANSLATION_CONTAINER_ENCRYPT_BYTES STRUCTS

type TR_ENCRYPT_BYTES_STATUS = string

type TrEncryptBytesMessage struct {
	TranslationContainerName string `json:"translation_container_name"`
	EncryptionKey            []byte `json:"enc_key"`
	CryptoType               string `json:"crypto_type"`
	Message                  []byte `json:"message"`
	IncludeUUID              bool   `json:"include_uuid"`
	Base64ReturnMessage      bool   `json:"base64_message"`
}

type TrEncryptBytesMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message []byte `json:"message"`
}

package translationstructs

// TRANSLATION_CONTAINER_MYTHIC_C2_TO_CUSTOM_MESSAGE_FORMAT STRUCTS

type TrMythicC2ToCustomMessageFormatMessage struct {
	TranslationContainerName string                 `json:"translation_container_name"`
	C2Name                   string                 `json:"c2_profile_name"`
	Message                  map[string]interface{} `json:"message"`
	UUID                     string                 `json:"uuid"`
	MythicEncrypts           bool                   `json:"mythic_encrypts"`
	CryptoKeys               []CryptoKeys           `json:"crypto_keys"`
}

type TrMythicC2ToCustomMessageFormatMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message []byte `json:"message"`
}

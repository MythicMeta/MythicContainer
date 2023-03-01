package translationstructs

// TRANSLATION_CONTAINER_CUSTOM_MESSAGE_TO_MYTHIC_C2_FORMAT STRUCTS

type TrCustomMessageToMythicC2FormatMessage struct {
	TranslationContainerName string       `json:"translation_container_name"`
	C2Name                   string       `json:"c2_profile_name"`
	Message                  []byte       `json:"message"`
	UUID                     string       `json:"uuid"`
	MythicEncrypts           bool         `json:"mythic_encrypts"`
	CryptoKeys               []CryptoKeys `json:"crypto_keys"`
}

type TrCustomMessageToMythicC2FormatMessageResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Message map[string]interface{} `json:"message"`
}

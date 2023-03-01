package translationstructs

// TRANSLATION_CONTAINER_GENERATE_ENCRYPTION_KEYS STRUCTS

type TrGenerateEncryptionKeysMessage struct {
	TranslationContainerName string `json:"translation_container_name"`
	C2Name                   string `json:"c2_profile_name"`
	CryptoParamValue         string `json:"value"`
	CryptoParamName          string `json:"name"`
}

type TrGenerateEncryptionKeysMessageResponse struct {
	Success       bool    `json:"success"`
	Error         string  `json:"error"`
	EncryptionKey *[]byte `json:"enc_key"`
	DecryptionKey *[]byte `json:"dec_key"`
}

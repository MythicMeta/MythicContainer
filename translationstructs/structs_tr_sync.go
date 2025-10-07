package translationstructs

import "github.com/MythicMeta/MythicContainer/utils/sharedStructs"

type TranslationContainer struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// SemVer is a specific semantic version tracker you can use for your payload type
	SemVer                        string                                                                                    `json:"semver"`
	Author                        string                                                                                    `json:"author"`
	TranslateCustomToMythicFormat TranslateCustomToMythicFormatFunction                                                     `json:"-"`
	TranslateMythicToCustomFormat TranslateMythicToCustomFormatFunction                                                     `json:"-"`
	GenerateEncryptionKeys        GenerateEncryptionKeysFunction                                                            `json:"-"`
	EncryptBytes                  EncryptBytesFunction                                                                      `json:"-"`
	DecryptBytes                  DecryptBytesFunction                                                                      `json:"-"`
	OnContainerStartFunction      func(sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse `json:"-"`
}

// TR_SYNC STRUCTS

type TrSyncMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type TrSyncMessage struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	Author           string `json:"author"`
	ContainerVersion string `json:"container_version"`
}

type TranslateCustomToMythicFormatFunction = func(input TrCustomMessageToMythicC2FormatMessage) TrCustomMessageToMythicC2FormatMessageResponse
type TranslateMythicToCustomFormatFunction = func(input TrMythicC2ToCustomMessageFormatMessage) TrMythicC2ToCustomMessageFormatMessageResponse
type GenerateEncryptionKeysFunction = func(input TrGenerateEncryptionKeysMessage) TrGenerateEncryptionKeysMessageResponse
type EncryptBytesFunction = func(input TrEncryptBytesMessage) TrEncryptBytesMessageResponse
type DecryptBytesFunction = func(input TrDecryptBytesMessage) TrDecryptBytesMessageResponse

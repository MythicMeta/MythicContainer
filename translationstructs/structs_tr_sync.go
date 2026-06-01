package translationstructs

import (
	"context"

	"github.com/MythicMeta/MythicContainer/utils/sharedStructs"
)

type TranslationContainer struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// SemVer is a specific semantic version tracker you can use for your payload type
	SemVer                        string                                                                                                     `json:"semver"`
	Author                        string                                                                                                     `json:"author"`
	TranslateCustomToMythicFormat TranslateCustomToMythicFormatFunction                                                                      `json:"-"`
	TranslateMythicToCustomFormat TranslateMythicToCustomFormatFunction                                                                      `json:"-"`
	GenerateEncryptionKeys        GenerateEncryptionKeysFunction                                                                             `json:"-"`
	EncryptBytes                  EncryptBytesFunction                                                                                       `json:"-"`
	DecryptBytes                  DecryptBytesFunction                                                                                       `json:"-"`
	OnContainerStartFunction      func(context.Context, sharedStructs.ContainerOnStartMessage) sharedStructs.ContainerOnStartMessageResponse `json:"-"`
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

type TranslateCustomToMythicFormatFunction = func(context.Context, TrCustomMessageToMythicC2FormatMessage) TrCustomMessageToMythicC2FormatMessageResponse
type TranslateMythicToCustomFormatFunction = func(context.Context, TrMythicC2ToCustomMessageFormatMessage) TrMythicC2ToCustomMessageFormatMessageResponse
type GenerateEncryptionKeysFunction = func(context.Context, TrGenerateEncryptionKeysMessage) TrGenerateEncryptionKeysMessageResponse
type EncryptBytesFunction = func(context.Context, TrEncryptBytesMessage) TrEncryptBytesMessageResponse
type DecryptBytesFunction = func(context.Context, TrDecryptBytesMessage) TrDecryptBytesMessageResponse

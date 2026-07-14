package mythicrpc

import "context"

import "github.com/MythicMeta/MythicContainer/utils/mythicutils"

type MythicRPCPayloadGetContentMessage struct {
	PayloadUUID string `json:"uuid"`
}

// Every mythicRPC function call must return a response that includes the following two values
type MythicRPCPayloadGetContentMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Content []byte `json:"content"`
}

func SendMythicRPCPayloadGetContent(ctx context.Context, input MythicRPCPayloadGetContentMessage) (*MythicRPCPayloadGetContentMessageResponse, error) {
	response := MythicRPCPayloadGetContentMessageResponse{}
	directFileToken, err := getDirectFileToken(ctx, input.PayloadUUID, "download")
	if err != nil {
		response.Error = err.Error()
		response.Success = false
		return &response, nil
	}
	if contents, err := mythicutils.GetFileFromMythic(ctx, input.PayloadUUID, directFileToken); err != nil {
		response.Error = err.Error()
		response.Success = false
	} else {
		response.Success = true
		response.Content = *contents
	}
	return &response, nil
}

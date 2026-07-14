package mythicrpc

import (
	"context"

	"github.com/MythicMeta/MythicContainer/rabbitmq"
)

type MythicRPCDirectFileTokenCreateMessage = rabbitmq.DirectFileTokenCreateMessage

type MythicRPCDirectFileTokenCreateMessageResponse = rabbitmq.DirectFileTokenCreateMessageResponse

func SendMythicRPCDirectFileTokenCreate(ctx context.Context, input MythicRPCDirectFileTokenCreateMessage) (*MythicRPCDirectFileTokenCreateMessageResponse, error) {
	return rabbitmq.SendDirectFileTokenCreate(ctx, input)
}

func getDirectFileToken(ctx context.Context, fileID string, action string) (string, error) {
	return rabbitmq.RequestDirectFileToken(ctx, fileID, action)
}

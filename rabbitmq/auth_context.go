package rabbitmq

import (
	"context"
	"fmt"

	"github.com/MythicMeta/MythicContainer/utils/authcontext"
	amqp "github.com/rabbitmq/amqp091-go"
)

const MythicAuthContextHeader = authcontext.MythicAuthContextHeader

func NewContextWithAuthToken(ctx context.Context, authToken string) context.Context {
	return authcontext.NewContextWithAuthToken(ctx, authToken)
}

func AuthTokenFromContext(ctx context.Context) (string, bool) {
	return authcontext.AuthTokenFromContext(ctx)
}

func HeadersFromContext(ctx context.Context) amqp.Table {
	if authToken, ok := AuthTokenFromContext(ctx); ok {
		return amqp.Table{MythicAuthContextHeader: authToken}
	}
	return nil
}

func AuthTokenFromHeaders(headers amqp.Table) (string, bool) {
	if headers == nil {
		return "", false
	}
	value, ok := headers[MythicAuthContextHeader]
	if !ok {
		return "", false
	}
	switch typedValue := value.(type) {
	case string:
		return typedValue, typedValue != ""
	case []byte:
		token := string(typedValue)
		return token, token != ""
	default:
		token := fmt.Sprintf("%v", typedValue)
		return token, token != ""
	}
}

func ContextFromHeaders(headers amqp.Table) context.Context {
	ctx := context.Background()
	if authToken, ok := AuthTokenFromHeaders(headers); ok {
		return NewContextWithAuthToken(ctx, authToken)
	}
	return ctx
}

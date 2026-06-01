package authcontext

import "context"

type authContextKey struct{}

const MythicAuthContextHeader = "mythic-auth-context"

func NewContextWithAuthToken(ctx context.Context, authToken string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if authToken == "" {
		return ctx
	}
	return context.WithValue(ctx, authContextKey{}, authToken)
}

func AuthTokenFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	value, ok := ctx.Value(authContextKey{}).(string)
	return value, ok && value != ""
}

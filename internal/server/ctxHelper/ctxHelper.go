package ctxHelper

import "context"

type ctxKey uint16

const (
	CtxKeyEmail ctxKey = iota
	CtxKeyIsAdmin
	CtxKeyRequestID
)

func IsEmailMatch(ctx context.Context, email string) bool {
	ctxEmail, ok := ctx.Value(CtxKeyEmail).(string)
	if !ok {
		return false
	}
	return ctxEmail == email
}

func IsAdmin(ctx context.Context) bool {
	ctxAdmin, ok := ctx.Value(CtxKeyIsAdmin).(bool)
	if !ok {
		return false
	}
	return ctxAdmin
}

func RequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(CtxKeyRequestID).(string)
	if !ok {
		return ""
	}
	return requestID
}

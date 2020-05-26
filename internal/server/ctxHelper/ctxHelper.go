package ctxHelper

import (
	"context"
	"github.com/balabanovds/void/internal/models"
)

type ctxKey uint16

const (
	CtxKeyEmail ctxKey = iota
	CtxKeyRole
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
	role, ok := ctx.Value(CtxKeyRole).(models.Role)
	if !ok {
		return false
	}
	return role == models.Admin
}

func IsManager(ctx context.Context) bool {
	role, ok := ctx.Value(CtxKeyRole).(models.Role)
	if !ok {
		return false
	}
	return role == models.Manager
}

func RequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(CtxKeyRequestID).(string)
	if !ok {
		return ""
	}
	return requestID
}

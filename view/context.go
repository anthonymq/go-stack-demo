package view

import (
	"context"
)

type contextKey string

const ContextAvatarUrlKey = contextKey("avatarUrl")

func GetAvatarUrl(ctx context.Context) string {
	if avatarUrl, ok := ctx.Value(ContextAvatarUrlKey).(string); ok {
		return avatarUrl
	}
	return "LOL"

}

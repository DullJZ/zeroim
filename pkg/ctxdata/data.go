package ctxdata

import (
	"context"
)

func GetUid(ctx context.Context) string {
	// 从jwt claims里获取UID
	if u, ok := ctx.Value(Identify).(string); ok {
		return u
	} else{
		return ""
	}
}
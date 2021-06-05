package logger

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type (
	traceIDContextKey struct{}
)

// NewTraceIDContext 创建跟踪ID上下文
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDContextKey{}, traceID)
}

// FromTraceIDContext 从上下文中获取跟踪ID
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return uuid.NewV4().String()
}

func NewTraceLogger(ctx context.Context, logger *zap.Logger) *zap.Logger {
	traceId := FromTraceIDContext(ctx)
	return logger.With(zap.String("trace_id", traceId))
}

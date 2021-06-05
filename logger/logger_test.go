package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	ctx := context.Background()
	ctx = NewTraceIDContext(ctx, "123123123")
	logger := NewLogger(WithDebugLevel(), WithConsoleEncoder())
	defer logger.Sync()
	logger = NewTraceLogger(ctx, logger)
	logger.Debug("debug logger")
	logger.Info("info logger")
	logger.Error("error logger")
	logger = logger.With(zap.String("with", "32q4324"))
	logger.Info("with test")
}

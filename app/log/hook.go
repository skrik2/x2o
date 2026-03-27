package log

import (
	"context"

	"go.uber.org/zap"
)

type Hook interface {
	Apply(ctx context.Context, msg string, fields ...zap.Field) []zap.Field
}

package sqlx

import (
	"context"
	"time"

	"github.com/kanthorlabs/common/logging"
	"gorm.io/gorm/logger"
)

func NewLogger(logger logging.Logger) logger.Interface {
	return &Logger{logger: logger}
}

type Logger struct {
	logger logging.Logger
}

func (logger *Logger) LogMode(logger.LogLevel) logger.Interface {
	return logger
}

func (logger *Logger) Info(ctx context.Context, msg string, args ...any) {
	logger.logger.Infow(msg, args...)
}
func (logger *Logger) Warn(ctx context.Context, msg string, args ...any) {
	logger.logger.Warnw(msg, args...)
}

func (logger *Logger) Error(ctx context.Context, msg string, args ...any) {
	logger.logger.Errorw(msg, args...)
}

func (logger *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)

	sql, rows := fc()
	if sql == ReadinessQuery || sql == LivenessQuery {
		return
	}
	args := []any{
		"rows", rows,
		"time", float64(elapsed.Nanoseconds()) / 1e6,
	}
	if err != nil {
		args = append(args, "error", err.Error())
	}

	logger.logger.Debugw(sql, args...)
}

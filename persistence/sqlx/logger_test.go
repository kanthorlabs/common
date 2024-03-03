package sqlx

import (
	"context"
	"testing"
	"time"

	mocklogging "github.com/kanthorlabs/common/mocks/logging"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm/logger"
)

func TestLogger(t *testing.T) {
	mocklogger := mocklogging.NewLogger(t)

	sqlogger := NewLogger(mocklogger)
	// nothing will happen
	sqlogger.LogMode(logger.Error)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	msg, args := logs()

	mocklogger.EXPECT().Infow(msg, args...).Once()
	sqlogger.Info(ctx, msg, args...)

	mocklogger.EXPECT().Warnw(msg, args...).Once()
	sqlogger.Warn(ctx, msg, args...)

	mocklogger.EXPECT().Errorw(msg, args...).Once()
	sqlogger.Error(ctx, msg, args...)

	// ignore readiness and liveness queries
	sqlogger.Trace(ctx, time.Now(), func() (string, int64) {
		return ReadinessQuery, 0
	}, nil)
	sqlogger.Trace(ctx, time.Now(), func() (string, int64) {
		return LivenessQuery, 0
	}, nil)

	sql := testdata.Fake.Lorem().Sentences(1)[0]
	rows := testdata.Fake.Int64Between(10, 1000)
	mocklogger.EXPECT().Debugw(sql, "rows", rows, "time", mock.AnythingOfType("float64"), "error", testdata.ErrGeneric.Error()).Once()

	sqlogger.Trace(ctx, time.Now(), func() (string, int64) {
		return sql, rows
	}, testdata.ErrGeneric)
}

func logs() (string, []any) {
	msg := testdata.Fake.Internet().StatusCodeWithMessage()

	var args []any
	count := testdata.Fake.IntBetween(10, 100)
	for i := 0; i < count; i++ {
		args = append(args, testdata.Fake.Emoji().EmojiCode())
	}

	return msg, args
}

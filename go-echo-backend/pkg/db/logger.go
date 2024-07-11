package db

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Logger struct {
	ZapLogger                 *logger.Logger
	LogLevel                  gormLogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func DefaultLogger(namspace string) *Logger {
	var dbLogger = &Logger{
		ZapLogger:                 logger.New(namspace),
		LogLevel:                  gormLogger.Info,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
	}
	gormLogger.Default = dbLogger
	return dbLogger
}
func (l Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return Logger{
		ZapLogger:                 l.ZapLogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormLogger.Info {
		return
	}
	l.logger().Sugar().Debugf(str, args...)
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormLogger.Warn {
		return
	}
	l.logger().Sugar().Warnf(str, args...)
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormLogger.Error {
		return
	}
	l.logger().Sugar().Errorf(str, args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}

	var elapsed = time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		sql = helper.EscapeQuotes(sql)
		l.logger().Error("trace", zap.Error(err), zap.Duration("sql_elapsed", elapsed), zap.Int64("sql_rows", rows), zap.String("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		sql = helper.EscapeQuotes(sql)
		l.logger().Warn("trace", zap.Duration("sql_elapsed", elapsed), zap.Int64("sql_rows", rows), zap.String("sql", sql))
	case l.LogLevel >= gormLogger.Info:
		sql, rows := fc()
		sql = helper.EscapeQuotes(sql)
		l.logger().Debug("trace", zap.Duration("sql_elapsed", elapsed), zap.Int64("sql_rows", rows), zap.String("sql", sql))
	}
}

func (l Logger) logger() *zap.Logger {
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)

		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, "gorm.io"):
		case strings.Contains(file, "queryfunc"):
		case strings.Contains(file, "query"):
		default:
			return l.ZapLogger.WithOptions(zap.AddCallerSkip(i - 1))
		}
	}
	return l.ZapLogger.Logger
}

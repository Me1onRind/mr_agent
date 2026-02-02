package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLoggerOptions configures the slog-backed GORM logger.
type GormLoggerOptions struct {
	Level                     gormlogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// GormLogger is a slog-backed implementation of gorm/logger.Interface.
type GormLogger struct {
	level                     gormlogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

func NewGormLogger(opts GormLoggerOptions) gormlogger.Interface {
	if opts.Level == 0 {
		opts.Level = gormlogger.Warn
	}
	if opts.SlowThreshold == 0 {
		opts.SlowThreshold = 200 * time.Millisecond
	}
	return &GormLogger{
		level:                     opts.Level,
		slowThreshold:             opts.SlowThreshold,
		ignoreRecordNotFoundError: opts.IgnoreRecordNotFoundError,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	next := *l
	next.level = level
	return &next
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.level >= gormlogger.Info {
		LoggerFromCtx(ctx).Info(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.level >= gormlogger.Warn {
		LoggerFromCtx(ctx).Warn(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.level >= gormlogger.Error {
		LoggerFromCtx(ctx).Error(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.level >= gormlogger.Error &&
		(!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		sql, rows := fc()
		LoggerFromCtx(ctx).Error(
			"gorm query error",
			slog.String("error", err.Error()),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	case l.slowThreshold != 0 && elapsed > l.slowThreshold && l.level >= gormlogger.Warn:
		sql, rows := fc()
		LoggerFromCtx(ctx).Warn(
			"gorm slow query",
			slog.Duration("elapsed", elapsed),
			slog.Duration("slow_threshold", l.slowThreshold),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	case l.level == gormlogger.Info:
		sql, rows := fc()
		LoggerFromCtx(ctx).Info(
			"gorm query",
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
			slog.String("sql", sql),
		)
	}
}

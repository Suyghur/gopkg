//@File     logger.go
//@Time     2024/8/6
//@Author   #Suyghur,

package gormx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/logger"
)

var _ logger.Interface = (*gormxLogger)(nil)

type (
	gormxLogger struct {
		logger.Config
	}
)

func NewLogger() logger.Interface {
	return &gormxLogger{
		Config: logger.Config{
			SlowThreshold:             defaultSlowThreshold,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		},
	}
}

func (l *gormxLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *gormxLogger) Info(ctx context.Context, format string, args ...any) {
	if l.LogLevel >= logger.Info {
		logx.WithContext(ctx).Infof(format, args...)
	}
}

func (l *gormxLogger) Warn(ctx context.Context, format string, args ...any) {
	if l.LogLevel >= logger.Warn {
		logx.WithContext(ctx).Infof(format, args...)
	}
}

func (l *gormxLogger) Error(ctx context.Context, format string, args ...any) {
	if l.LogLevel >= logger.Error {
		logx.WithContext(ctx).Errorf(format, args...)
	}
}

func (l *gormxLogger) Trace(ctx context.Context, begin time.Time, fn func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fn()
	//// 通用字段
	logFields := []logx.LogField{
		logx.Field("sql", sql),
		logx.Field("duration", microsecondsStr(elapsed)),
		logx.Field("rows", rows),
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		if rows == -1 {
			logx.WithContext(ctx).Infow("[GORMX-WARN]: exec sql on record not found", logFields...)
		} else {
			logx.WithContext(ctx).Errorw(fmt.Sprintf("[GORMX-ERROR]: exec sql on error: %v", err), logFields...)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		logx.WithContext(ctx).Sloww("[GORMX-WARN]: exec sql on slow", logFields...)
	case l.LogLevel == logger.Info:
		logx.WithContext(ctx).Infow("[GORMX-INFO]: exec sql", logFields...)
	}
}

// microsecondsStr 将 time.Duration 类型（nano seconds 为单位）输出为小数点后 3 位的 ms （microsecond 毫秒，千分之一秒）
func microsecondsStr(elapsed time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
}

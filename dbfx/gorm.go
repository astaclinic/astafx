// modified from default logger for zap
// https://github.com/go-gorm/gorm/blob/master/logger/logger.go

package dbfx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

var ErrRecordNotFound = errors.New("record not found")

type GormLoggerConfig struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

type GormLogger struct {
	logger *zap.SugaredLogger
	config GormLoggerConfig
}

func NewGormLogger(logger *zap.SugaredLogger) *GormLogger {
	return &GormLogger{logger, GormLoggerConfig{
		SlowThreshold:             time.Second, // Slow SQL threshold
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
	}}
}

func (g *GormLogger) LogMode(logger.LogLevel) logger.Interface {
	return g
}

func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	g.logger.Info(msg, data)
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	g.logger.Warn(msg, data)
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	g.logger.Error(msg, data)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	msg := "Executed SQL statement"
	switch {
	case err != nil && (!errors.Is(err, ErrRecordNotFound) || !g.config.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			g.logger.Errorw(msg, "err", err, "time", float64(elapsed.Nanoseconds())/1e6, "rows", "-", "sql", sql)
		} else {
			g.logger.Errorw(msg, "err", err, "time", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case elapsed > g.config.SlowThreshold && g.config.SlowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.config.SlowThreshold)
		if rows == -1 {
			g.logger.Warnw(msg, "err", slowLog, "time", float64(elapsed.Nanoseconds())/1e6, "rows", "-", "sql", sql)
		} else {
			g.logger.Warnw(msg, "err", slowLog, "time", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			g.logger.Infow(msg, "time", float64(elapsed.Nanoseconds())/1e6, "rows", "-", "sql", sql)
		} else {
			g.logger.Infow(msg, "time", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	}
}

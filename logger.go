package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/atom-providers/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Logger struct {
	Level logger.LogLevel
}

func (g *Logger) LogMode(l logger.LogLevel) logger.Interface {
	g.Level = l
	return g
}

func (g *Logger) Info(_ context.Context, msg string, params ...interface{}) {
	log.Infof(msg, params...)
}

func (g *Logger) Warn(_ context.Context, msg string, params ...interface{}) {
	log.Warnf(msg, params)
}

func (g *Logger) Debug(_ context.Context, msg string, params ...interface{}) {
	log.Debugf(msg, params)
}

func (g *Logger) Error(_ context.Context, msg string, params ...interface{}) {
	log.Errorf(msg, params)
}

func (g *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	var (
		traceStr     = "%s [%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s[%.3fms] [rows:%v] %s"
	)

	elapsed := time.Since(begin)
	SlowThreshold := time.Second

	sql, rows := fc()
	switch {
	case err != nil && g.Level >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		if rows == -1 {
			g.Error(ctx, traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.Error(ctx, traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > SlowThreshold && SlowThreshold != 0 && g.Level >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", SlowThreshold)
		if rows == -1 {
			g.Warn(ctx, traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.Warn(ctx, traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case g.Level == logger.Info:
		if rows == -1 {
			g.Debug(ctx, traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.Debug(ctx, traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

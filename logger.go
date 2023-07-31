package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rogeecn/atom-addons/providers/log"
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
	SlowThreshold := time.Second * 2

	sql, rows := fc()
	switch {
	case err != nil && g.Level >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		if rows == -1 {
			log.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > SlowThreshold && SlowThreshold != 0 && g.Level >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", SlowThreshold)
		if rows == -1 {
			log.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case g.Level == logger.Info:
		if rows == -1 {
			log.Debugf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Debugf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

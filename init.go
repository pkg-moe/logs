package logs

import (
	"context"

	"go.uber.org/zap"

	"pkg.moe/pkg/logger"
)

type Field = zap.Field

type LogsLevel int64

const (
	Error LogsLevel = iota
	Info
	Debug
)

type Logs struct {
	log *zap.Logger
	ctx context.Context
}

type logsItem struct {
	log   *zap.Logger
	msg   string
	level LogsLevel
	field []Field
}

var ctxFunc func(ctx context.Context) *zap.Logger = logger.GetWithContext

func InitLogsContextProvider(ctxFuncProvider func(ctx context.Context) *zap.Logger) {
	ctxFunc = ctxFuncProvider
}

func NewLogs(ctx context.Context) *Logs {
	return &Logs{ctx: ctx, log: ctxFunc(ctx)}
}

func (l *Logs) Ctx() context.Context {
	return l.ctx
}

func (l *Logs) Error(msg string, err error) *logsItem {
	log := &logsItem{l.log, msg, Error, []Field{}}
	log.field = append(log.field, zap.Error(err))
	return log
}

func (l *Logs) Info(msg string) *logsItem {
	return &logsItem{l.log, msg, Info, []Field{}}
}

func (l *Logs) Debug(msg string) *logsItem {
	return &logsItem{l.log, msg, Debug, []Field{}}
}

func (l *logsItem) Tag(key string, value interface{}) *logsItem {
	l.field = append(l.field, zap.Any(key, value))
	return l
}

func (l *logsItem) Apply() {
	log := l.log.WithOptions(zap.AddCallerSkip(1))
	switch l.level {
	case Error:
		log.Error(l.msg, l.field...)
	case Info:
		log.Info(l.msg, l.field...)
	case Debug:
		log.Debug(l.msg, l.field...)
	}
}

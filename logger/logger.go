package logger

import "go.uber.org/zap"

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	DPanic(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Close()
}

type logger struct {
	base *zap.SugaredLogger
}

func NewLogger() (Logger, error) {
	lgr, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	sugar := lgr.Sugar()
	l := &logger{
		base: sugar,
	}

	return l, nil
}

func (l *logger) Debug(msg string, args ...interface{}) {
	l.base.Debugw(msg, args)
}

func (l *logger) Info(msg string, args ...interface{}) {
	l.base.Infow(msg, args)
}

func (l *logger) Warn(msg string, args ...interface{}) {
	l.base.Warnw(msg, args)
}

func (l *logger) Error(msg string, args ...interface{}) {
	l.base.Errorw(msg, args)
}

func (l *logger) DPanic(msg string, args ...interface{}) {
	l.base.DPanicw(msg, args)
}

func (l *logger) Panic(msg string, args ...interface{}) {
	l.base.Panicw(msg, args)
}

func (l *logger) Fatal(msg string, args ...interface{}) {
	l.base.Fatalw(msg, args)
}

func (l *logger) Close() {
	l.base.Sync()
}

package logger

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/misikdmitriy/password-sharing/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerFactory interface {
	NewLogger() (*zap.Logger, func(), error)
}

type loggerFactory struct {
	configuration *config.Config
}

type testLoggerFactory struct {
}

func NewLoggerFactory(configuration *config.Config) LoggerFactory {
	return &loggerFactory{
		configuration: configuration,
	}
}

func NewTestLoggerFactory() LoggerFactory {
	return &testLoggerFactory{}
}

func (lf *loggerFactory) NewLogger() (*zap.Logger, func(), error) {
	now := time.Now()
	logfile := path.Join(lf.configuration.Zap.LogsPath, fmt.Sprintf("%s.log", now.Format("2006-01-02-15")))

	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), lf.configuration.Zap.Level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), lf.configuration.Zap.Level),
	)

	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	close := func() {
		log.Sync()
		file.Close()
	}

	return log, close, nil
}

func (lf *testLoggerFactory) NewLogger() (*zap.Logger, func(), error) {
	return zap.NewExample(), func() {}, nil
}

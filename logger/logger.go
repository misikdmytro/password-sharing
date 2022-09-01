package logger

import (
	"os"

	"github.com/misikdmitriy/password-sharing/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CloseFunc func()

func NewLogger(c *config.Config) (*zap.Logger, CloseFunc, error) {
	file, err := os.Create(c.Zap.LogsPath)
	if err != nil {
		return nil, nil, err
	}

	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), c.Zap.Level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), c.Zap.Level),
	)

	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))

	return log, func() {
		log.Sync()
		file.Close()
	}, nil
}

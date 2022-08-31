package logger

import (
	"os"

	"github.com/misikdmitriy/password-sharing/config"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

func NewLogger(c *config.Config) *zap.Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, c.Zap.Level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	return logger
}

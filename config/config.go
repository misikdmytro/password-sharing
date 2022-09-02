package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Database struct {
		ConnectionString string `mapstructure:"connectionstring"`
		Provider         string `mapstructure:"provider"`
	} `mapstructure:"database"`
	App struct {
		LinkLength    int    `mapstructure:"linklength"`
		Port          int    `mapstructure:"port"`
		Address       string `mapstructure:"address"`
		ConsulAddress string `mapstructure:"consuladdress"`
		ServiceId     int    `mapstructure:"serviceid"`
		BasePath      string `mapstructure:"basepath"`
	} `mapstructure:"app"`
	Zap struct {
		Level    zapcore.Level `mapstructure:"level"`
		LogsPath string        `mapstructure:"logspath"`
	} `mapstructure:"zap"`
	Encrypt struct {
		Secret string `mapstructure:"secret"`
		IV     []byte `mapstructure:"iv"`
	} `mapstructure:"encrypt"`
}

func LoadConfig() (*Config, error) {
	conf := viper.New()

	configname := fmt.Sprintf("%s.yaml", getWebEnv())
	conf.SetConfigFile(configname)
	conf.SetConfigType("yaml")

	conf.SetEnvPrefix("psconfig")
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	conf.AutomaticEnv()

	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	c := &Config{}
	err := conf.Unmarshal(c)

	return c, err
}

const defaultEnv = "dev"

func getWebEnv() string {
	env := os.Getenv("WEB_ENV")
	if env != "" {
		return env
	}

	return defaultEnv
}

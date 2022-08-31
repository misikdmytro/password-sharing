package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Database struct {
		ConnectionString string `mapstructure:"connectionstring"`
		MaxConnection    int    `mapstructure:"maxconnection"`
		Timeout          int    `mapstructure:"timeout"`
		Provider         string `mapstructure:"provider"`
	} `mapstructure:"database"`
	App struct {
		LinkLength int `mapstructure:"linklength"`
		Port       int `mapstructure:"port"`
	} `mapstructure:"app"`
	Zap struct {
		Level zapcore.Level `mapstructure:"level"`
	} `mapstructure:"zap"`
}

func CreateEmpty() *Config {
	return &Config{}
}

func LoadConfig() (*Config, error) {
	conf := viper.New()

	conf.SetConfigFile(getWebEnv() + ".json")
	conf.SetConfigType("json")

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

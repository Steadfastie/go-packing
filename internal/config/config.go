package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func Load() (Config, error) {
	env := strings.TrimSpace(os.Getenv("APP_ENV"))
	if env == "" {
		env = "dev"
	}
	if env != "dev" && env != "prod" {
		return Config{}, fmt.Errorf("unsupported APP_ENV %q, expected dev or prod", env)
	}

	v := viper.New()
	v.SetConfigType("json")
	v.SetConfigName(env)
	v.AddConfigPath("./internal/config")

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read %s config: %w", env, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

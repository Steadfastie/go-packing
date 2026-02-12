package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv     string        `mapstructure:"app_env"`
	Server     ServerConfig   `mapstructure:"server"`
	Database   DatabaseConfig `mapstructure:"database"`
	Log        LogConfig      `mapstructure:"log"`
	SourcePath string        `mapstructure:"-"`
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
	v.AddConfigPath("./cmd/config")
	v.AddConfigPath("/app/cmd/config")
	v.SetDefault("log.level", "info")

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read %s config: %w", env, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}
	if strings.TrimSpace(cfg.Database.URL) == "" {
		return Config{}, fmt.Errorf("database.url is required (from config file or DATABASE_URL)")
	}
	if strings.TrimSpace(cfg.Server.Port) == "" {
		return Config{}, fmt.Errorf("server.port is required (from config file or PORT)")
	}
	cfg.SourcePath = v.ConfigFileUsed()

	return cfg, nil
}

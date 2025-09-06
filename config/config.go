package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
}

func Parse() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config") // config.yaml
	v.SetConfigType("yaml")
	v.AddConfigPath("./config") // путь к yaml

	// 1) читаем YAML (если не найден — логируем как error, но продолжаем)
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "config yaml error")
	}

	// 2) ENV перекрывают YAML
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "config : unmarshal failed")
	}
	return &cfg, nil
}

func (c *Config) PGURL() string {
	return "postgres://" + c.DB.User + ":" + c.DB.Password +
		"@" + c.DB.Host + ":" + c.DB.Port + "/" + c.DB.Name +
		"?sslmode=" + c.DB.SSLMode
}

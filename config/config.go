package config

import (
	"fmt"

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
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	// 1)reading from YAML (if not found - looging like an err with continuing)
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "config yaml error")
	}

	// 2) ENV taking over YAML
	v.AutomaticEnv()
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "config : unmarshal failed")
	}
	return &cfg, nil
}

func (c *Config) PGURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
		c.DB.SSLMode,
	)
}

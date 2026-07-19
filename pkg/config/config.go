// Package config loads typed application configuration from environment
// variables.
package config

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

// Config gives typed access to configuration values by key.
type Config interface {
	GetString(key string) string
	GetInt(key string) int
}

type config struct {
	v *viper.Viper
}

func New() (Config, error) {
	v := viper.New()

	v.SetDefault("http.addr", ":8080")
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "console")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return &config{v: v}, nil
}

func (c *config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *config) GetInt(key string) int {
	return c.v.GetInt(key)
}

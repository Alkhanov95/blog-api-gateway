package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mtvy/blog-api-gateway/internal/logger"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//go:embed config.yml
var defaultYamlFile []byte

type Config struct {
	HTTP struct {
		Port int
	}
	Log logger.Config
}

func Parse() (*Config, error) {
	return parse(defaultYamlFile, "BLOG_APIGATEWAY")
}

func (c *Config) GetHTTPEndpoint() string {
	return fmt.Sprintf(":%d", c.HTTP.Port)
}

func parse(defaultYamlFile []byte, envPrefix string) (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadConfig(bytes.NewBuffer(defaultYamlFile)); err != nil {
		return nil, errors.Wrap(err, "reading default settings")
	}

	if err := viper.MergeInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	} else {
		if file := viper.ConfigFileUsed(); file != "" {
			slog.Info(fmt.Sprintf("configuration file: %s", file))
		}
	}

	if _, err := os.Stat("./.env"); !errors.Is(err, os.ErrNotExist) {
		if err = godotenv.Load("./.env"); err != nil {
			return nil, errors.Wrap(err, "read .env file")
		} else {
			slog.Info("using .env file")
		}
	}

	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "parse settings")
	}

	return cfg, nil
}

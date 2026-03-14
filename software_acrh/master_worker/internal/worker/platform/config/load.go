package config

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigFile(".env")
	v.SetConfigType("env")

	if _, err := os.Stat(".env"); err == nil {
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
	} else {
		slog.Warn("config file not found, using environment variables")
	}

	// allow env override
	v.AutomaticEnv()
	autoBindEnv(v, Config{})

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Errorf("agent load configuration error: %w", err))
	}

	return cfg
}

func autoBindEnv(v *viper.Viper, s any) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("mapstructure")
		if tag != "" {
			_ = v.BindEnv(tag)
		}
	}
}

func validate(cfg Config) error {
	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

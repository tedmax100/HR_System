package config

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	DataBaseURL   string `mapstructure:"DATABASE_URL"`
	RedisURL      string `mapstructure:"REDIS_URL"`
	Env           string `mapstructure:"ENV"`
	Port          int    `mapstructure:"PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiration int    `mapstructure:"JWT_EXPIRATION"`
}

func NewConfig() (cfg *Config, err error) {
	cfg = &Config{}
	loader := viper.New()

	if err = cfg.bindEnvsRecursive(loader, reflect.ValueOf(cfg)); err != nil {
		return nil, err
	}

	loader.AutomaticEnv()
	if err := loader.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	return cfg, nil

}

func (cfg *Config) bindEnvsRecursive(loader *viper.Viper, val reflect.Value) error {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i)

		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			tag = field.Name
		}

		switch value.Kind() {
		case reflect.Struct:
			if err := cfg.bindEnvsRecursive(loader, value); err != nil {
				return err
			}
		default:
			if err := loader.BindEnv(tag); err != nil {
				return fmt.Errorf("error binding env var: %s: %w", tag, err)
			}
		}
	}

	return nil
}

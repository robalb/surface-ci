// Package envconfig contains primitive functions for parsing the program environment variables into
// a configuration struct.
// This is not the main program configuration. The only paramenters expected via ENV variables
// are API secrets. Everything else will be read from configuration files. See: pkg/configfiles
package envconfig

import (
	"log/slog"
	"reflect"
	"strings"
)

type EnvConfig struct {
	ConfigFolder string `env:"CONFIG_FOLDER"`
	DataFolder   string `env:"DATA_FOLDER"`
	SecretTest   string `env:"SECRET_TEST" sensitive:"true"`
}

func defaultEnvConfig() EnvConfig {
	return EnvConfig{
		ConfigFolder: "./",
		DataFolder:   "./data",
		SecretTest:   "",
	}
}

func New(
	args []string,
	getenv func(string) string,
	log *slog.Logger,
) (EnvConfig, error) {

	c := defaultEnvConfig()

	v := reflect.ValueOf(&c).Elem()
	t := reflect.TypeOf(c)

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		envName := fieldType.Tag.Get("env")
		if envName == "" {
			continue
		}

		envValue := getenv(envName)
		if envValue == "" {
			continue
		}

		field.SetString(envValue)

		// censor the value before logging it
		isSensitive := fieldType.Tag.Get("sensitive") == "true"
		if isSensitive {
			envValue = strings.Repeat("*", len(envValue))
		}

		//log default variable overrides
		log.Info("Read ENV variable", "name", envName, "value", envValue)
	}

	return c, nil
}

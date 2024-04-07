package config

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Read(configFile string, out any) error {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	hooks := mapstructure.ComposeDecodeHookFunc(
		DurationHook(),
	)

	err = viper.Unmarshal(out, viper.DecodeHook(hooks))
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ProgramName       = "gifseekr"
	DefaultConfigName = "settings"
)

type Settings struct {
	GiphyAPIKey string `mapstructure:"giphy-api-key"`
}

func ReadConfig() (settings *Settings, err error) {
	v := viper.New()

	// define options
	flagset := pflag.NewFlagSet(ProgramName, pflag.ExitOnError)
	flagset.StringP("config-dir", "D", "", "path to config directory")
	flagset.StringP("config", "C", DefaultConfigName, "name of config file, without extension")
	flagset.String("giphy-api-key", "", "a Giphy API key")

	// read command line args
	if err = flagset.Parse(os.Args); err != nil {
		return nil, err
	}
	if err = v.BindPFlags(flagset); err != nil {
		return nil, err
	}

	// read environment vars
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// read config file
	if configDir := v.GetString("config-dir"); configDir != "" {
		v.AddConfigPath(configDir)
	}
	v.AddConfigPath(".")
	configName := v.GetString("config")
	v.SetConfigName(configName)
	err = v.MergeInConfig()
	// config file is optional so long as it was not set by the user
	if _, ok := err.(viper.ConfigFileNotFoundError); ok && configName != DefaultConfigName {
		err = nil
	}

	settings = new(Settings)
	err = v.Unmarshal(settings)
	return settings, err
}

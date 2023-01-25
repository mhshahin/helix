package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Cfg Config

// Config contains all the required configs and
// environment variables of the project.
type Config struct {
	Namespace string `yaml:"namespace"`
	Mediums   Medium `yaml:"mediums"`
	Rules     []Rule `yaml:"rules"`
}

type Medium struct {
	Matrix Matrix `yaml:"matrix"`
}

// MatrixConfig contains required data to post messages to Matrix
// server's webook and publish to a channel
type Matrix struct {
	DisplayName string `yaml:"display_name"`
	Address     string `yaml:"address"`
	Token       string `yaml:"token"`
	Timeout     string `yaml:"timeout"`
}

type Rule struct {
	Kind  string `yaml:"kind"`
	Label string `yaml:"label"`
	Type  string `yaml:"type"`
}

func LoadConfig(configPath string) {
	viper.SetEnvPrefix("helix")
	viper.SetConfigName(configPath)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&Cfg)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}

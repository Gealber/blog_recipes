package config

import (
	"github.com/spf13/viper"
)

const (
	DevEnvironment = "DEV"
)

type AppConfig struct {
	TON struct {
		Net    string
		Wallet struct {
			Seed []string
		}
	}
}

var (
	cfg *AppConfig
)

func Config() *AppConfig {
	if cfg == nil {
		loadConfig()
	}

	return cfg
}

func loadConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Ignore config file not found, perhaps we will use environment variables.
	_ = viper.ReadInConfig()

	cfg = &AppConfig{}

	// TON.
	cfg.TON.Wallet.Seed = viper.GetStringSlice("TON_WALLET_SEED")
}

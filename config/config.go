package config

import "github.com/spf13/viper"

type Config struct {
	// Define your configuration variables here
	Mongo_URI string
	JWT_TOKEN string
}

func SetUpConfig() (*Config, error) {
	// Use -> VIPER <- to read the environement variables
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	config := &Config{
		Mongo_URI: viper.GetString("Mongo_URI"),
		JWT_TOKEN: viper.GetString("JWT_TOKEN"),
	}
	return config, nil
}

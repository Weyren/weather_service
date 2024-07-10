package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"database"`
	API struct {
		Key      string `mapstructure:"key"`
		Interval int    `mapstructure:"interval"`
	} `mapstructure:"api"`
}

// LoadConfig loads config file from path and returns Config struct or error
func LoadConfig(path string) (*Config, error) {
	var config *Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	//Read in config file
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	//Unmarshal config file
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}

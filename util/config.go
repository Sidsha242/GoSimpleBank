package util

import "github.com/spf13/viper"

//will store all configuration for the application
//the values are read by viper from a config file or rnv file

type Config struct {
	DBDriber string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}


//LoadConfig reads configuration from file or env variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path) //location of congif file
	viper.SetConfigName("app") //look for config file with name app
	viper.SetConfigType("env") //type of config file also json, yaml etc

	viper.AutomaticEnv() //bring values from environment variables automatically

	err = viper.ReadInConfig() //read the config file
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return //return config object we can use this in main.go
}
package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Server struct {
		Port string `yml:"port"`
		Host string `yml:"host"`
	} `yml:"server"`
	DB struct {
		Host     string `yml:"host"`
		Port     int    `yml:"port"`
		User     string `yml:"user"`
		Password string `yml:"password"`
		DBName   string `yml:"dbname"`
		SSLMode  string `yml:"sslmode"`
	} `yml:"db"`
}

// LoadConfig ...
func MustLoad() *Config{
	viper.SetConfigName("config")  
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs") 

	// to use existing enviroment virables 
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`,`_`))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Cannot read config: %v", err)
		
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling: %v", err)
		
	}

	return &config
}

package config

import (
	"falconapi/internal/model"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func InitViper() *model.Config {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.SetEnvPrefix("demo")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("unable to initialize viper: %w", err))
	}

	var config model.Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to unmarshal config: %s", err)
	}

	log.Println("viper config initialized")

	return &config
}

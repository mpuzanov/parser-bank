package config

import (
	"log"
	"net"
	"strings"

	"github.com/mpuzanov/parser-bank/pkg/logger"
	"github.com/spf13/viper"
)

//Config Структура фйла с конфигурацией
type Config struct {
	Log      logger.LogConf `yaml:"log" mapstructure:"log"`
	HTTPAddr string         `yaml:"http_listen" mapstructure:"http_listn"`
	Host     string         `yaml:"http_host" mapstructure:"http_host"`
	Port     string         `yaml:"http_port" mapstructure:"http_port"`
}

// LoadConfig Загрузка конфигурации из файла
func LoadConfig(filePath string) (*Config, error) {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("log.level", "info")
	viper.SetDefault("http_host", "0.0.0.0")
	viper.SetDefault("http_port", "7777")

	if filePath != "" {
		log.Printf("Parsing config: %s\n", filePath)
		viper.SetConfigFile(filePath)
		viper.SetConfigType("yaml")
		//log.Println(viper.ConfigFileUsed())
		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("Config file is not specified.")
	}
	//log.Println(viper.AllSettings())

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	config.HTTPAddr = net.JoinHostPort(config.Host, config.Port)
	//log.Println(config)
	return &config, nil
}

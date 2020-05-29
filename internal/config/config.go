package config

import (
	"net"
	"strings"

	"github.com/mpuzanov/parser-bank/pkg/logger"
	"github.com/spf13/viper"
)

//Config Структура фйла с конфигурацией
type Config struct {
	Log      logger.LogConf `yaml:"log" mapstructure:"log"`
	HTTPAddr string         `yaml:"http_listen" mapstructure:"http_listn"`
	Host     string         `yaml:"host" mapstructure:"host"`
	Port     string         `yaml:"port" mapstructure:"port"`
	PathTmp  string         `yaml:"path_tmp" mapstructure:"path_tmp"`
}

// LoadConfig Загрузка конфигурации из файла
func LoadConfig(filePath string) (*Config, error) {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format_JSON", "true")
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", "7777")
	viper.SetDefault("path_tmp", "./tmp_files/")

	if filePath != "" {
		logger.LogSugar.Debugf("Parsing config: %s", filePath)
		viper.SetConfigFile(filePath)
		viper.SetConfigType("yaml")
		//logger.LogSugar.Debug(viper.ConfigFileUsed())
		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	} else {
		logger.LogSugar.Debug("Config file is not specified.")
	}
	//logger.LogSugar.Debug(viper.AllSettings())

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	config.HTTPAddr = net.JoinHostPort(config.Host, config.Port)
	//logger.LogSugar.Debug(config)
	return &config, nil
}

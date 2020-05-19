package main

import (
	"github.com/mpuzanov/parser-bank/internal/config"
	"github.com/mpuzanov/parser-bank/internal/web"
	"github.com/mpuzanov/parser-bank/pkg/logger"
	flag "github.com/spf13/pflag"
)

func main() {
	var cfgPath string
	flag.StringVarP(&cfgPath, "config", "c", "", "path to the configuration file")
	flag.Parse()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		logger.LogSugar.Fatalf("Не удалось загрузить %s: %s", cfgPath, err)
	}
	l := logger.NewLogger(cfg.Log)

	if err := web.Start(cfg, l); err != nil {
		logger.LogSugar.Fatal(err)
	}
}

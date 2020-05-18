package main

import (
	"log"

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
		log.Fatalf("Не удалось загрузить %s: %s", cfgPath, err)
	}
	logger := logger.NewLogger(cfg.Log)

	if err := web.Start(cfg, logger); err != nil {
		log.Fatal(err)
	}
}

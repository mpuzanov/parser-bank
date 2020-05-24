package web

import (
	"github.com/mpuzanov/parser-bank/internal/config"
	"github.com/mpuzanov/parser-bank/internal/web"
	"github.com/mpuzanov/parser-bank/pkg/logger"
	"github.com/spf13/cobra"
)

var cfgPath string

func init() {
	ServerCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "path to the configuration file")
}

var (
	// ServerCmd .
	ServerCmd = &cobra.Command{
		Use:   "web_server",
		Short: "Run web server",
		Run:   webServerStart,
	}
)

func webServerStart(cmd *cobra.Command, args []string) {

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		logger.LogSugar.Fatalf("Не удалось загрузить %s: %s", cfgPath, err)
	}
	l := logger.NewLogger(cfg.Log)

	if err := web.Start(cfg, l); err != nil {
		logger.LogSugar.Fatal(err)
	}
}

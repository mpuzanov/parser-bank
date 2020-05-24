package main

import (
	"github.com/mpuzanov/parser-bank/pkg/logger"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.LogSugar.Fatal(err)
	}
}

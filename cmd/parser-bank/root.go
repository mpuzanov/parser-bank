package main

import (
	"github.com/mpuzanov/parser-bank/cmd/parser-bank/shell"
	"github.com/mpuzanov/parser-bank/cmd/parser-bank/web"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "parser-bank",
	Short: "",
}

func init() {
	rootCmd.AddCommand(web.ServerCmd, shell.RunCmd)
}

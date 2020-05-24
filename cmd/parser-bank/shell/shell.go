package shell

import (
	"os"
	"path/filepath"

	"github.com/mpuzanov/parser-bank/pkg/logger"

	"github.com/spf13/cobra"
)

var pathFiles string

var (
	// RunCmd .
	RunCmd = &cobra.Command{
		Use:   "shell",
		Short: "Run shell",
		Run:   pathStart,
	}
)

func init() {
	RunCmd.Flags().StringVarP(&pathFiles, "path", "p", "", "path to process")
}

func pathStart(cmd *cobra.Command, args []string) {
	logger.LogSugar.Info("Обработка: ", pathFiles)
	//Посетить все файлы и папки в дереве каталогов
	err := filepath.Walk(pathFiles,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			logger.LogSugar.Infof("%s - %d", path, info.Size())
			return nil
		})
	if err != nil {
		logger.LogSugar.Error(err)
	}
}

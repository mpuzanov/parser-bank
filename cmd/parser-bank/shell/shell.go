package shell

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mpuzanov/parser-bank/internal/domain/model"
	"github.com/mpuzanov/parser-bank/internal/storage"
	"github.com/mpuzanov/parser-bank/internal/storage/payments"
	"github.com/mpuzanov/parser-bank/pkg/logger"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

type result struct {
	fileName string
	val      []model.Payment
	err      error
}

// type storageFormat struct {
// 	db *storage.ListFormatBanks
// }

type storageFormat struct {
	fb *storage.ListFormatBanks
}

var (
	pathFiles string
	debug     bool
	log       *zap.Logger
	store     storageFormat //storage.ListFormatBanks

	// RunCmd .
	RunCmd = &cobra.Command{
		Use:   "shell",
		Short: "Run shell",
		Run:   pathStart,
	}
)

func init() {
	RunCmd.Flags().StringVarP(&pathFiles, "path", "p", "", "path to process")
	RunCmd.Flags().BoolVarP(&debug, "debug", "d", false, "run debug")
}

func pathStart(cmd *cobra.Command, args []string) {
	level := "info"
	if debug {
		level = "debug"
	}
	log = logger.InitLogger(logger.LogConf{Level: level})
	zap.ReplaceGlobals(log)

	files := []string{}
	log.Sugar().Info("Обработка: ", pathFiles)
	//Посетить все файлы и папки в дереве каталогов
	err := filepath.Walk(pathFiles,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			//log.Sugar().Debugf("%s - %d", path, info.Size())
			files = append(files, path)
			return nil
		})
	if err != nil {
		log.Sugar().Error(err)
	}

	countFiles := len(files)
	if countFiles == 0 {
		log.Sugar().Info("Файлов не выбрано")
		return
	}
	log.Sugar().Infof("Файлов для обработки: %d", countFiles)
	countWorker := getCountWorker(countFiles)

	store = storageFormat{fb: storage.NewFormatBanks()}
	if err := store.fb.Open(); err != nil {
		log.Sugar().Fatalf("error load format banks %v", err)
	}
	log.Sugar().Debugf("%v", store.fb)
	chanFile := make(chan string)    // канал обработки файлов
	chanResults := make(chan result) // канал получения результатов
	chanStop := make(chan struct{})  // канал для прерывания выполнения горутин

	startProg := time.Now()

	wg := sync.WaitGroup{}
	// запускаем обработчики
	for i := 0; i < countWorker; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			store.worker(i, chanFile, chanResults, chanStop)
		}(i)
	}
	log.Sugar().Infof("Запущенно обработчиков: %d", countWorker)

	go func() { // ждём окончания всех обработчиков
		defer close(chanResults)
		wg.Wait()
	}()

	go func() { // отправляем файлы на обработку
		defer close(chanFile)
		for i := 0; i < countFiles; i++ {
			chanFile <- files[i]
		}
	}()

	count := 0
	countError := 0
	valuesTotal := payments.ListPayments{}
	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(1)
	// получаем результаты
	go func() {
		defer wgReceivers.Done()
		// Читаем из канала результатов
		for r := range chanResults {
			count++
			if r.err != nil {
				countError++
				log.Sugar().Errorf("Файл: %s (%v)", r.fileName, r.err)
				continue
			}
			valuesTotal.Db = append(valuesTotal.Db, r.val...)
		}
	}()
	wgReceivers.Wait()
	log.Sugar().Infof("Обработано файлов: %d, с ошибками: %d", count, countError)
	log.Sugar().Infof("Итого платежей: %d", len(valuesTotal.Db))
	log.Sugar().Info("Время выполнения: ", time.Since(startProg))

}

func (s *storageFormat) worker(id int, files <-chan string, res chan<- result, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			//log.Sugar().Debugf("Останавливаем обработчик %d", id)
			return
		default:
		}

		select {
		case <-stop:
			//log.Sugar().Debugf("Останавливаем обработчик %d", id)
			return
		case file, ok := <-files:
			if !ok {
				//log.Sugar().Debugf("Заданий больше нет - Останавливаем обработчик %d", id)
				return
			}
			r := result{fileName: file}
			r.val, r.err = s.fb.ReadFile(file, log) //log много выдаёт информации надо упорядочить  zap.L()
			res <- r
			//log.Sugar().Debugf("обработчик %d выполнил задание! Ошибка: %v", id, r.err)
		}

	}
}

// getCountWorker определяем кол-во обработчиков для запуска
func getCountWorker(countFiles int) int {
	count := countFiles / 3
	if count == 0 {
		count = 1
	}
	if count > 10 {
		count = 10
	}
	return count
}

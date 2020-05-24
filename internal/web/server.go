package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mpuzanov/parser-bank/internal/config"
	"github.com/mpuzanov/parser-bank/internal/repository"
	"github.com/mpuzanov/parser-bank/internal/storage"

	"go.uber.org/zap"
)

type myHandler struct {
	router *mux.Router
	logger *zap.Logger
	cfg    *config.Config
	store  repository.StorageFormatBanks
}

// Start .
func Start(conf *config.Config, log *zap.Logger) error {

	err := preparePath(conf.PathTmp)
	if err != nil {
		log.Error("Prepare path", zap.Error(err))
		os.Exit(1)
	}

	handler := &myHandler{
		router: mux.NewRouter(),
		logger: log,
		cfg:    conf,
		store:  storage.NewFormatBanks(),
	}
	handler.configRouter()

	handler.loadStore()

	srv := &http.Server{
		Addr:           conf.HTTPAddr,
		Handler:        handler,
		IdleTimeout:    10 * time.Second,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//запускаем веб-сервер
	go func(log *zap.Logger) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen", zap.Error(err))
		}
	}(log)
	log.Info("Starting Http server", zap.String("address", srv.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	//подчищаем каталог в временными файлами
	err = deleteTmpFiles(conf.PathTmp)
	if err != nil {
		log.Error("deleteTmpFiles", zap.Error(err))
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed", zap.Error(err))
	}
	log.Info("Shutdown done")
	os.Exit(0)

	return nil

}

//preparePath проверяем существуют ли необходимые каталоги для работы
func preparePath(tmpDir string) error {
	//проверяем существует ли каталог временных файлов
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		err := os.Mkdir(tmpDir, 0666)
		if err != nil {
			return err
		}
	}
	//проверяем существует ли каталог для загрузки файлов
	dir := path.Join(tmpDir, "in")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0666)
		if err != nil {
			return err
		}
	}
	//проверяем существует ли каталог для выгрузки файлов
	dir = path.Join(tmpDir, "out")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteTmpFiles(tmpDir string) error {

	dir := path.Join(tmpDir, "in")
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir, 0666)
	if err != nil {
		return err
	}

	dir = path.Join(tmpDir, "out")
	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir, 0666)
	if err != nil {
		return err
	}
	return nil
}

package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
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
	store  repository.StorageFormatBanks
}

// Start .
func Start(conf *config.Config, log *zap.Logger) error {

	handler := &myHandler{
		router: mux.NewRouter(),
		logger: log,
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
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed", zap.Error(err))
	}
	log.Info("Shutdown done")
	os.Exit(0)

	return nil

}

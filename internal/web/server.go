package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mpuzanov/parser-bank/internal/config"
	"github.com/mpuzanov/parser-bank/internal/store"
	"go.uber.org/zap"
)

type myHandler struct {
	router *mux.Router
	logger *zap.Logger
	store  *store.Store
}

// Start .
func Start(conf *config.Config, logger *zap.Logger) error {

	handler := &myHandler{
		router: mux.NewRouter(),
		logger: logger,
		store: store.New()  ,
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
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Server http started: %s, file log: %s\n", srv.Addr, conf.Log.File)
	logger.Info("Starting Http server", zap.String("address", srv.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Print("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}
	log.Println("Shutdown done")
	os.Exit(0)

	return nil

}

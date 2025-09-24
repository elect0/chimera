package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/elect0/chimera/internal/config"
	"github.com/elect0/chimera/internal/logger"
)

func main() {
	cfg := config.New()

	log := logger.New(cfg.Log.Level)
	log.Info("logger initialized", slog.String("level", cfg.Log.Level))
	log.Info("starting chimera service")

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpSever.Port),
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("starting http server", slog.Int("port", cfg.HttpSever.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	sig := <-quit
	log.Info("received shutdown signal", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HttpSever.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("server shutdown gracefully")
}

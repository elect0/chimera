package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/elect0/chimera/internal/adapters/api"
	"github.com/elect0/chimera/internal/adapters/cache"
	"github.com/elect0/chimera/internal/adapters/storage"
	"github.com/elect0/chimera/internal/application/transformation"
	"github.com/elect0/chimera/internal/config"
	"github.com/elect0/chimera/internal/logger"
)

func main() {
	cfg := config.New()

	log := logger.New(cfg.Log.Level)
	log.Info("logger initialized", slog.String("level", cfg.Log.Level))
	fmt.Println(`
          _             _       _     _         _   _         _            _           _          
        /\ \           / /\    / /\  /\ \      /\_\/\_\ _    /\ \         /\ \        / /\        
       /  \ \         / / /   / / /  \ \ \    / / / / //\_\ /  \ \       /  \ \      / /  \       
      / /\ \ \       / /_/   / / /   /\ \_\  /\ \/ \ \/ / // /\ \ \     / /\ \ \    / / /\ \      
     / / /\ \ \     / /\ \__/ / /   / /\/_/ /  \____\__/ // / /\ \_\   / / /\ \_\  / / /\ \ \     
    / / /  \ \_\   / /\ \___\/ /   / / /   / /\/________// /_/_ \/_/  / / /_/ / / / / /  \ \ \    
   / / /    \/_/  / / /\/___/ /   / / /   / / /\/_// / // /____/\    / / /__\/ / / / /___/ /\ \   
  / / /          / / /   / / /   / / /   / / /    / / // /\____\/   / / /_____/ / / /_____/ /\ \  
 / / /________  / / /   / / /___/ / /__ / / /    / / // / /______  / / /\ \ \  / /_________/\ \ \ 
/ / /_________\/ / /   / / //\__\/_/___\\/_/    / / // / /_______\/ / /  \ \ \/ / /_       __\ \_\
\/____________/\/_/    \/_/ \/_________/        \/_/ \/__________/\/_/    \_\/\_\___\     /____/_/
                                                                                                  
		`)

	originRepo, err := storage.NewS3OriginRepository(context.Background(), cfg, log)
	if err != nil {
		log.Error("failed to create S3 origin repository", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("S3 origin repository initialized")

	cacheRepo, err := cache.NewRedisCacheRepository(context.Background(), cfg, log)
	if err != nil {
		log.Error("failed to create redis cache repository", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("redis cache repository initialized")

	transformationService := transformation.NewService(log, originRepo, cacheRepo)

	apiHandler := api.NewHandler(transformationService, log)

	mux := http.NewServeMux()

	apiHandler.RegisterRoutes(mux)

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

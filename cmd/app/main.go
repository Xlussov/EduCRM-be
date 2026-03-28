package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Xlussov/EduCRM-be/docs"
	"github.com/Xlussov/EduCRM-be/internal/app"
	"github.com/Xlussov/EduCRM-be/pkg/config"
	"github.com/Xlussov/EduCRM-be/pkg/logger"
)

// @title EduCRM API
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := app.New(ctx, cfg, log)
	if err != nil {
		panic("the application crashed due to a critical error")
	}

	a.Start(ctx)

	<-ctx.Done()
	log.Info("shutdown signal received, stopping...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Stop(shutdownCtx); err != nil {
		log.Errorf("graceful shutdown error: %v", err)
	}

	log.Info("app stopped gracefully")
}

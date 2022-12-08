package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"

	"github.com/astaclinic/astafx/logger"
)

func Run(mainApp *fx.App) {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Graceful shutdown
	// https://github.com/gin-gonic/gin#manually

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := mainApp.Start(ctx); err != nil {
			logger.Fatalf("Error in staring application: %v", err)
			os.Exit(1)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	logger.Infof("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mainApp.Stop(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown.")
		os.Exit(1)
	}

	logger.Infof("Server exiting")
}

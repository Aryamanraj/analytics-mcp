package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/payram/payram-analytics-mcp-server/internal/agent/admin"
	"github.com/payram/payram-analytics-mcp-server/internal/agent/supervisor"
	"github.com/payram/payram-analytics-mcp-server/internal/logging"
)

func main() {
	addr := os.Getenv("PAYRAM_AGENT_LISTEN_ADDR")
	if addr == "" {
		addr = ":9900"
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	sup, err := supervisor.NewFromEnv()
	if err != nil {
		log.Fatalf("failed to configure supervisor: %v", err)
	}
	if err := sup.Start(ctx); err != nil {
		log.Fatalf("failed to start supervisor: %v", err)
	}

	handler := admin.NewMux(sup)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	logger, cleanup, err := logging.New("agent")
	useLogger := err == nil
	if useLogger {
		defer cleanup()
		logger.Infof("agent starting on %s", addr)
	} else {
		log.Printf("agent starting on %s (fallback logging): %v", addr, err)
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if useLogger {
				logger.Errorf("server error: %v", err)
			} else {
				log.Printf("server error: %v", err)
			}
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil && err != http.ErrServerClosed {
		if useLogger {
			logger.Errorf("shutdown error: %v", err)
		} else {
			log.Printf("shutdown error: %v", err)
		}
	}

	sup.Wait()
}

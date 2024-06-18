package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pchaivong/pong-service/internal/health"
	"github.com/pchaivong/pong-service/internal/pong"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	fmt.Println("Starting up....")

	// Create context that listen for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Get PORT environment variable
	// Format PORT = ':8080'
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	router := gin.Default()
	// Register healthcheck
	hcHandler := health.NewHandler()
	router.GET("/healthz/status", hcHandler.HealthCheck)

	// Register prometheus metrics
	router.GET("/metrics", func(c *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Writer, c.Request)
	})
	// Register handlers here
	svcHandler := pong.NewHandler()
	router.GET("/v1/pong", svcHandler.Pong)
	// Create http server
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Initializing the server in the goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Printf("Starting server, listen port = %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Starting server error: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

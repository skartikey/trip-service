package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tripservice/api/route"
)

func main() {
	loadEnv()
	serveApplication()
}

func loadEnv() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func serveApplication() {
	router := gin.Default()
	route.Routing(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Channel to listen for OS signals for shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		// Context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
		log.Println("Server exiting")
	}()

	log.Println("Starting server on port 8080...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

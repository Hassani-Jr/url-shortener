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

	"github.com/Hassani-Jr/url-shortener/internal/handler"
	"github.com/Hassani-Jr/url-shortener/internal/middleware"
	"github.com/Hassani-Jr/url-shortener/internal/service"
	"github.com/Hassani-Jr/url-shortener/internal/storage"
	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(); err != nil{
		log.Println("No .env file found, using environment variables")
	}

	port := getEnv("SERVER_PORT", "8080")
	// Initiallize dependencies
	store := storage.NewMemoryStorage()
	svc := service.NewShortenerService(store)
	urlHandler := handler.NewURLHandler(svc)




	// 1. Create http server
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", urlHandler.Shorten)
	mux.HandleFunc("GET /{code}", urlHandler.Redirect)
	mux.HandleFunc("GET /urls/{code}/stats", urlHandler.Stats)
	mux.HandleFunc("DELETE /{code}",urlHandler.Delete)
	mux.HandleFunc("GET /health",healthHandler)

	handler := middleware.RequestID(loggingMiddleware(mux))

 
	server := &http.Server{
		Addr: ":"+ port,
		Handler: handler,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	// start server in goroutine
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal,1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
	ctx,cancel := context.WithTimeout(context.Background(),30 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx);err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped gracefully")
}

func healthHandler (w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
	}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start := time.Now()
		requestID := middleware.GetRequestID(r.Context())
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("[%s] Completed in %v", requestID ,time.Since(start))
	})
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

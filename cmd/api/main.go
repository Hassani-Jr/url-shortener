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
    "github.com/Hassani-Jr/url-shortener/internal/service"
    "github.com/Hassani-Jr/url-shortener/internal/storage"
)

func main(){
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

 
	server := &http.Server{
		Addr: ":8080",
		Handler: loggingMiddleware(mux),
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
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

func uniqueRequest


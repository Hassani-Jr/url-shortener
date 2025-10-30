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
)

func main(){
	// 1. Create http server
	mux := http.NewServeMux()
	mux.HandleFunc("/health",healthHandler)


	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
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

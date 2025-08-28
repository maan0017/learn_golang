package main

import (
	"context"
	"go/web-sockets/internal/chat"
	"go/web-sockets/internal/httpx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	addr := getEnv("ADDR", ":8080")
	allowedOrigins := getEnv("WS_ALLOWED_ORIGINS", "http://localhost:3000")

	hub := chat.NewHub()
	go hub.Run()

	// HTTP server + handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.Handle("/ws", httpx.WebSocketHandler(hub, allowedOrigins))

	srv := &http.Server{
		Addr:              addr,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run server
	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	hub.Close()
	log.Println("bye 👋")
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

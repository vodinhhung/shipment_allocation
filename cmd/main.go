package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shipment_allocation/internal/api"
	"shipment_allocation/internal/dependency"
	"time"

	"shipment_allocation/cmd/router"
	"shipment_allocation/internal/common"
)

func main() {
	//rand.Seed(time.Now().UnixNano())

	// config via env with fallbacks
	addr := common.GetEnv("ADDR", ":8080")
	readTimeout := common.ParseDurationOrDefault(common.GetEnv("READ_TIMEOUT", "5s"), 5*time.Second)
	writeTimeout := common.ParseDurationOrDefault(common.GetEnv("WRITE_TIMEOUT", "10s"), 10*time.Second)
	idleTimeout := common.ParseDurationOrDefault(common.GetEnv("IDLE_TIMEOUT", "120s"), 120*time.Second)

	db, err := dependency.InitDb()
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	allocation := api.NewAllocation(db)
	rout := router.NewRouter(allocation)

	mux := http.NewServeMux()
	mux.HandleFunc("/", rout.HandleRoot)
	mux.HandleFunc("/health", rout.HandleHealth)
	mux.HandleFunc("/allocate", rout.HandleAllocateShipment)

	handler := simpleCORSMiddleware(loggingMiddleware(mux))
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		// listen for interrupt signal
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("shutdown: signal received, shutting down HTTP server...")
		// allow up to 10s for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("starting server on %s", addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}

	<-idleConnsClosed
	log.Println("server stopped")
}

// loggingMiddleware logs basic request info and duration.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := fmt.Sprintf("%d", start.UnixNano())
		// attach req-id to response for tracing
		w.Header().Set("X-Request-ID", reqID)
		log.Printf("[req=%s] %s %s from=%s", reqID, r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("[req=%s] done in %s", reqID, time.Since(start))
	})
}

// simpleCORS adds permissive CORS header for demo purposes.
func simpleCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

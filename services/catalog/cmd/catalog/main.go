package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jx-nin/sretail/services/catalog/internal/product"
	"github.com/jx-nin/sretail/services/catalog/internal/router"
)

func main() {
	store := product.NewMemoryStore([]product.Product{
		{
			ID:         "prod-1",
			Name:       "Mechanical Keyboard",
			PriceCents: 12999,
		},
		{
			ID:         "prod-2",
			Name:       "Ergonomic Mouse",
			PriceCents: 10999,
		},
	})

	s := &http.Server{
		Addr:              ":8081",
		Handler:           router.New(store),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MiB
	}

	serverErrors := make(chan error, 1)

	shutdownSignal, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		serverErrors <- s.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}

	case <-shutdownSignal.Done():
		stop()
		shutdownContext, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)

		defer cancel()
		if err := s.Shutdown(shutdownContext); err != nil {
			log.Printf("graceful shutdown failed: %v", err)

			if closeErr := s.Close(); closeErr != nil {
				log.Printf("force close failed: %v", closeErr)
			}
		}
	}
}

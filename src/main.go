package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/silviorneto/go-crud-in-memory/api"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to execute code: ", "error", err)
		os.Exit(1)
	}

	slog.Info("System offline")
}

func run() error {
	handler := api.NewHandler()

	s := http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

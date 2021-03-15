package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	auth "github.com/menucha-de/Util.Auth/auth"
	loglib "github.com/menucha-de/logging"
)

var lg *loglib.Logger = loglib.GetLogger("http")

func main() {
	var port = flag.Int("p", 80, "port")
	flag.Parse()

	router := auth.NewRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	errs := make(chan error)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- err
		}
	}()

	go auth.Init()

	lg.Infof("HTTP server started on port %d", *port)

	select {
	case err := <-errs:
		lg.Errorf("Failed to start HTTP server: %s", err)
	case <-done:
		lg.Info("Server stopped")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	auth.Shutdown()

	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		lg.Error("Server shutdown failed:", err.Error())
	}

	lg.Info("Server exited properly")
}

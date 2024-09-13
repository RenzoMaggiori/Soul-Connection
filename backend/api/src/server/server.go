package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"soul-connection.com/api/src/lib"
)

func Start(server *http.Server, running chan struct{}) {
	running <- struct{}{}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}
}

func Stop(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	lib.ServerLog("INFO", "Shutting down server...")

	if err := server.Close(); err != nil {
		log.Fatal("Server shutdown failed:", err)
	}
	lib.ServerLog("INFO", "Server gracefully stopped")
}

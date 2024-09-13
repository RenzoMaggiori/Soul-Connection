package server

import (
	"fmt"
	"net/http"
	"syscall"
	"testing"
)

func TestStartAndStop(t *testing.T) {
	apiServer := &http.Server{
		Addr: fmt.Sprintf(":%d", 8080),
	}
	running := make(chan struct{})
	go Start(apiServer, running)
	<-running
	go func() {
		defer Stop(apiServer)
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/pravinkanna/jQueue/internal/server"
)

const (
	port            = 50051
	shutdownTimeout = 15 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	log.Println("Server Shutdown Gracefully")
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("TCP server failed to start: %w", err)
	}

	s := grpc.NewServer()
	server.Register(s)
	reflection.Register(s)

	serveErr := make(chan error, 1)
	log.Printf("Server starting on port %d", port)
	go func() {
		serveErr <- s.Serve(lis)
	}()

	select {
	case err := <-serveErr:
		if err != nil {
			return fmt.Errorf("failed to serve: %w", err)
		}
	case <-ctx.Done():
		return shutdown(s)
	}

	return nil
}

// shutdown drains s, forcing termination if timeout elapses
func shutdown(s *grpc.Server) error {
	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
		// drained cleanly — return immediately
		return nil
	case <-time.After(shutdownTimeout):
		// took too long — force it
		s.Stop()
		return fmt.Errorf("graceful shutdown timed out after %s", shutdownTimeout)
	}

}

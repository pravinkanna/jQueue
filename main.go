package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/pravinkanna/jQueue/internal/server"
)

const port = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln("TCP server failed to start", err)
	}

	s := grpc.NewServer()
	server.Register(s)
	reflection.Register(s)

	log.Printf("Server starting on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/pravinkanna/jQueue/gen/go/jqueue/v1"
)

const port = 50051

type HealthServer struct {
	pb.UnimplementedHealthServiceServer
}

func (hs *HealthServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Message: "Alive"}, nil
}

type jobServer struct {
	pb.UnimplementedJobServiceServer
}

type queueServer struct {
	pb.UnimplementedQueueServiceServer
}

type leaseServer struct {
	pb.UnimplementedLeaseServiceServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln("TCP server failed to start", err)
	}

	s := grpc.NewServer()
	pb.RegisterHealthServiceServer(s, &HealthServer{})
	pb.RegisterJobServiceServer(s, &jobServer{})
	pb.RegisterQueueServiceServer(s, &queueServer{})
	pb.RegisterLeaseServiceServer(s, &leaseServer{})

	reflection.Register(s)

	log.Printf("Server starting on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

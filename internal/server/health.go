package server

import (
	"context"

	pb "github.com/pravinkanna/jQueue/gen/go/jqueue/v1"
)

type healthServer struct {
	pb.UnimplementedHealthServiceServer
}

func (hs *healthServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Message: "Alive"}, nil
}

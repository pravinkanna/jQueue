package server

import (
	"google.golang.org/grpc"

	pb "github.com/pravinkanna/jQueue/gen/go/jqueue/v1"

	"github.com/pravinkanna/jQueue/internal/store"
	"github.com/pravinkanna/jQueue/internal/store/memory"
)

var st store.Store = &memory.Memory{}

// Register attaches the gRPC services to gRPC server
func Register(s *grpc.Server) {
	pb.RegisterHealthServiceServer(s, &healthServer{})
	pb.RegisterJobServiceServer(s, &jobServer{})
	pb.RegisterQueueServiceServer(s, &queueServer{})
	pb.RegisterLeaseServiceServer(s, &leaseServer{})
}

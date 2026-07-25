package server

import (
	"google.golang.org/grpc"

	pb "github.com/pravinkanna/jQueue/gen/go/jqueue/v1"

	"github.com/pravinkanna/jQueue/internal/store"
)

// Register attaches the gRPC services to gRPC server
func Register(s *grpc.Server, st store.Store) {
	pb.RegisterHealthServiceServer(s, &healthServer{})
	pb.RegisterJobServiceServer(s, &jobServer{st: st})
	pb.RegisterQueueServiceServer(s, &queueServer{st: st})
	pb.RegisterLeaseServiceServer(s, &leaseServer{st: st})
}

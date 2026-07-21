package server

import (
	pb "github.com/pravinkanna/jQueue/gen/go/jqueue/v1"
)

type leaseServer struct {
	pb.UnimplementedLeaseServiceServer
}

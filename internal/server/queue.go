package server

import (
	"context"

	pb "github.com/pravinkanna/jQueue/gen/go/jqueue/v1"
	"github.com/pravinkanna/jQueue/internal/store"
)

type queueServer struct {
	pb.UnimplementedQueueServiceServer
	st store.Store
}

func (qs *queueServer) CreateQueue(ctx context.Context, req *pb.CreateQueueRequest) (*pb.CreateQueueResponse, error) {
	queueName := req.QueueName

	err := qs.st.CreateQueue(ctx, queueName)
	if err != nil {
		return nil, err
	}
	res := &pb.CreateQueueResponse{}
	return res, nil
}

func (qs *queueServer) DeleteQueue(ctx context.Context, req *pb.DeleteQueueRequest) (*pb.DeleteQueueResponse, error) {
	queueName := req.QueueName
	err := qs.st.DeleteQueue(ctx, queueName)
	if err != nil {
		return nil, err
	}
	res := &pb.DeleteQueueResponse{}
	return res, nil
}

func (qs *queueServer) PurgeQueue(ctx context.Context, req *pb.PurgeQueueRequest) (*pb.PurgeQueueResponse, error) {
	queueName := req.QueueName
	purgedCount, err := qs.st.PurgeQueue(ctx, queueName)
	if err != nil {
		return nil, err
	}
	res := &pb.PurgeQueueResponse{
		PurgedCount: purgedCount,
	}
	return res, nil
}

func (qs *queueServer) GetQueueStatus(ctx context.Context, req *pb.GetQueueStatusRequest) (*pb.GetQueueStatusResponse, error) {
	queueName := req.QueueName
	queue, err := qs.st.GetQueueStatus(ctx, queueName)
	if err != nil {
		return nil, err
	}
	res := &pb.GetQueueStatusResponse{
		Queue: &pb.Queue{
			Name:           queue.Name,
			PendingCount:   queue.PendingCount,
			ScheduledCount: queue.ScheduledCount,
			LeasedCount:    queue.LeasedCount,
			CompletedCount: queue.CompletedCount,
			DlqCount:       queue.DLQCount,
		},
	}
	return res, nil
}

func (qs *queueServer) ListQueues(ctx context.Context, req *pb.ListQueuesRequest) (*pb.ListQueuesResponse, error) {
	queues, err := qs.st.ListQueues(ctx)
	if err != nil {
		return nil, err
	}

	pbQueues := []*pb.Queue{}
	for _, item := range queues {
		q := &pb.Queue{
			Name:           item.Name,
			PendingCount:   item.PendingCount,
			ScheduledCount: item.ScheduledCount,
			LeasedCount:    item.LeasedCount,
			CompletedCount: item.CompletedCount,
			DlqCount:       item.DLQCount,
		}

		pbQueues = append(pbQueues, q)
	}

	res := &pb.ListQueuesResponse{
		Queues: pbQueues,
	}
	return res, nil
}

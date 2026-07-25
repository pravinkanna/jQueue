package server

import (
	"context"

	pb "github.com/pravinkanna/jQueue/gen/go/jqueue/v1"
	"github.com/pravinkanna/jQueue/internal/store"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type leaseServer struct {
	pb.UnimplementedLeaseServiceServer
	st store.Store
}

func (ls *leaseServer) LeaseJobs(ctx context.Context, req *pb.LeaseJobsRequest) (*pb.LeaseJobsResponse, error) {
	queueName := req.Queue
	batchSize := req.BatchSize
	leaseDuration := req.LeaseDuration.AsDuration()
	leasedJobs, err := ls.st.LeaseJobs(ctx, queueName, batchSize, leaseDuration)
	if err != nil {
		return nil, err
	}
	pbLeasedJobs := []*pb.LeasedJob{}
	for _, j := range leasedJobs {
		pbLeasedJob := &pb.LeasedJob{
			LeaseToken: j.LeaseToken,
			Job: &pb.Job{
				JobId:          j.Job.JobID,
				IdempotencyKey: j.Job.IdempotencyKey,
				Queue:          j.Job.Queue,
				Payload:        j.Job.Payload,
				State:          pb.JobState(j.Job.State),
				MaxRetries:     j.Job.MaxRetries,
				RetryCount:     j.Job.RetryCount,
				LastError:      j.Job.LastError,
				CreatedAt:      timestamppb.New(j.Job.CreatedAt),
				ScheduledAt:    timestamppb.New(j.Job.ScheduledAt),
				CompletedAt:    timestamppb.New(j.Job.CompletedAt),
			},
			LeaseExpiresAt: timestamppb.New(j.ExpiresAt),
		}
		pbLeasedJobs = append(pbLeasedJobs, pbLeasedJob)
	}

	res := &pb.LeaseJobsResponse{
		LeasedJobs: pbLeasedJobs,
	}
	return res, nil
}

func (ls *leaseServer) ExtendJobLease(ctx context.Context, req *pb.ExtendJobLeaseRequest) (*pb.ExtendJobLeaseResponse, error) {
	leaseToken := req.LeaseToken
	duration := req.Duration.AsDuration()
	expiresAt, err := ls.st.ExtendJobLease(ctx, leaseToken, duration)
	if err != nil {
		return nil, err
	}
	res := &pb.ExtendJobLeaseResponse{
		LeaseExpiresAt: timestamppb.New(expiresAt),
	}
	return res, nil
}

func (ls *leaseServer) AckJob(ctx context.Context, req *pb.AckJobRequest) (*pb.AckJobResponse, error) {
	leaseToken := req.LeaseToken
	err := ls.st.AckJob(ctx, leaseToken)
	if err != nil {
		return nil, err
	}
	res := &pb.AckJobResponse{}
	return res, nil
}

func (ls *leaseServer) NackJob(ctx context.Context, req *pb.NackJobRequest) (*pb.NackJobResponse, error) {
	leaseToken := req.LeaseToken
	reason := req.Reason
	err := ls.st.NackJob(ctx, leaseToken, reason)
	if err != nil {
		return nil, err
	}
	res := &pb.NackJobResponse{}
	return res, nil
}

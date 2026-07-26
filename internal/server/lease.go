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

func (ls *leaseServer) LeaseJobs(ctx context.Context, req *pb.LeaseJobRequest) (*pb.LeaseJobResponse, error) {
	queueName := req.Queue
	leaseDuration := req.LeaseDuration.AsDuration()
	leasedJobs, err := ls.st.LeaseJob(ctx, queueName, leaseDuration)
	if err != nil {
		return nil, err
	}
	pbLeasedJob := &pb.LeasedJob{
		LeaseToken: leasedJobs.LeaseToken,
		Job: &pb.Job{
			JobId:          leasedJobs.Job.JobID,
			IdempotencyKey: leasedJobs.Job.IdempotencyKey,
			Queue:          leasedJobs.Job.Queue,
			Payload:        leasedJobs.Job.Payload,
			State:          pb.JobState(leasedJobs.Job.State),
			MaxRetries:     leasedJobs.Job.MaxRetries,
			RetryCount:     leasedJobs.Job.RetryCount,
			LastError:      leasedJobs.Job.LastError,
			CreatedAt:      timestamppb.New(leasedJobs.Job.CreatedAt),
			ScheduledAt:    timestamppb.New(leasedJobs.Job.ScheduledAt),
			CompletedAt:    timestamppb.New(leasedJobs.Job.CompletedAt),
		},
		LeaseExpiresAt: timestamppb.New(leasedJobs.ExpiresAt),
	}

	res := &pb.LeaseJobResponse{
		LeasedJob: pbLeasedJob,
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

package server

import (
	"context"
	"fmt"
	"time"

	pb "github.com/pravinkanna/jQueue/gen/go/jqueue/v1"
	"github.com/pravinkanna/jQueue/internal/store"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type jobServer struct {
	pb.UnimplementedJobServiceServer
	st store.Store
}

func (js *jobServer) EnqueueJob(ctx context.Context, req *pb.EnqueueJobRequest) (*pb.EnqueueJobResponse, error) {
	runAt := req.GetRunAt()
	delay := req.GetDelay()
	runAtTs := time.Now()
	if runAt != nil && delay == nil {
		runAtTs = runAt.AsTime()
	} else {
		runAtTs = runAtTs.Add(delay.AsDuration())
	}

	params := store.EnqueueParams{
		IdempotencyKey: req.IdempotencyKey,
		Queue:          req.Queue,
		Payload:        req.Payload,
		MaxRetries:     req.MaxRetries,
		RunAt:          runAtTs,
	}
	jobID, isDup, err := js.st.EnqueueJob(ctx, params)
	if err != nil {
		return nil, err
	}
	res := &pb.EnqueueJobResponse{
		JobId:       jobID,
		IsDuplicate: isDup,
	}
	return res, nil
}

func (js *jobServer) GetJobState(ctx context.Context, req *pb.GetJobStateRequest) (*pb.GetJobStateResponse, error) {
	jobID := req.JobId
	job, err := js.st.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	res := &pb.GetJobStateResponse{
		State:       pb.JobState(job.State),
		RetryCount:  job.RetryCount,
		LastError:   job.LastError,
		ScheduledAt: timestamppb.New(job.ScheduledAt),
	}
	return res, nil
}

func (js *jobServer) CancelJob(ctx context.Context, req *pb.CancelJobRequest) (*pb.CancelJobResponse, error) {
	jobID := req.JobId
	state, err := js.st.CancelJob(ctx, jobID)
	if err != nil {
		return nil, err
	}
	res := &pb.CancelJobResponse{
		State: pb.JobState(state),
	}
	return res, nil
}

func (js *jobServer) ListJobs(ctx context.Context, req *pb.ListJobsRequest) (*pb.ListJobsResponse, error) {
	queueName := req.Queue
	state := store.JobState(req.StateFilter)
	pSize := req.PageSize
	pToken := req.PageToken
	jobs, nToken, err := js.st.ListJobs(ctx, queueName, state, pSize, pToken)
	if err != nil {
		return nil, err
	}

	var pbJobs []*pb.Job
	for _, j := range jobs {
		pbJob := &pb.Job{
			JobId:          j.JobID,
			IdempotencyKey: j.IdempotencyKey,
			Queue:          j.Queue,
			Payload:        j.Payload,
			State:          pb.JobState(j.State),
			MaxRetries:     j.MaxRetries,
			RetryCount:     j.RetryCount,
			LastError:      j.LastError,
			CreatedAt:      timestamppb.New(j.CreatedAt),
			ScheduledAt:    timestamppb.New(j.ScheduledAt),
			CompletedAt:    timestamppb.New(j.CompletedAt),
		}

		pbJobs = append(pbJobs, pbJob)
	}
	res := &pb.ListJobsResponse{
		Jobs:          pbJobs,
		NextPageToken: nToken,
	}
	return res, nil
}

func (js *jobServer) DLQRetryJobs(ctx context.Context, req *pb.DLQRetryJobsRequest) (*pb.DLQRetryJobsResponse, error) {
	jobID := req.GetJobId()
	queueName := req.GetQueue()

	var retryCount uint32
	var err error
	if jobID != "" {
		err = js.st.RetryDLQJob(ctx, jobID)
	} else if queueName != "" {
		retryCount, err = js.st.RetryDLQQueue(ctx, queueName)
	} else {
		return nil, fmt.Errorf("Provide either JobID or QueueName")
	}
	if err != nil {
		return nil, err
	}
	res := &pb.DLQRetryJobsResponse{
		RetriedCount: retryCount,
	}
	return res, nil
}

# Distributed Job Queue

## Requirements

### Functional Requirements
1. Clients can create jobs and workers can claim jobs and execute 
2. Client can schedule jobs immediately or later
3. Clients can get status of the job
4. Workers can listen to a queue in the server (long polling -> unary)
5. Workers can fetch one job or multiple jobs (batch)
6. Workers can Ack/Nack a job after the execution of a Job
7. CLI should be able to do all actions of clients, workers plus Queue Creation and Queue Status
8. User can access our prometheus with our metrics endpoint for observability

### Non Functional Requirements
1. Server should persist job on disk
2. Server should ensure at least once delivery
3. Server should be highly available
4. If no response from a worker (Ack/Nack/ExtendLease), consider it failed and retry (exponential backoff + jitter)
5. Graceful shutdown - finish current job before exit
6. Periodic sweeper will monitor the failed, timedout jobs back to either pending or DLQ (depends on retry limit)
7. The completed jobs will be cleaned up after 24 hours of completion
8. Server should deduplicate the jobs with idempotent key to prevent double creation of jobs

## Core Entities
1. Jobs -> States: Scheduled, Pending, Leased, Completed, Failed, DLQ, Cancelled
2. Schedule
3. Queue
4. Client
5. Workers

## APIs
1. rpc Enqueue()
2. rpc Status()
3. rpc Lease() -> returns leaseToken
4. rpc Ack(leaseToken)
5. rpc Nack(leaseToken)
6. rpc ExtendLease(leaseToken)
7. rpc Cancel()
8. rpc ListJobs()
9. rpc CreateQueue()
10. rpc DeleteQueue()
11. rpc QueueStatus()
12. rpc ListQueues()
13. rpc DLQRetry()

## Data Flow
Client ---Create Job---> Server(Broker) <---Lease Job--- Worker
Client <---Get Status---> Server(Broker) <---Ack/Nack/ExtendLease--- Worker

## Implementation Details

1. gRPC as API for both client and workers
2. Retry limit set by client
3. DLQ for failures exceeding retry limit
4. No ordering guarantee

## Plan
Phase 1:
No persistence - All data in memory

Phase 2:
Persist data - No need for distributed server (broker)
Shared DB (some SQL DB)

Phase 3:
should support distributed server (broker) support
Shared DB (some SQL DB - But should handle concurrent queries (ATOMIC))

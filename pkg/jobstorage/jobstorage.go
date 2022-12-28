package jobstorage

import "context"

// BlobJobStorage is an interface for jobs dispatch.
type JobStorage interface {
	UpdateBlobList(ctx context.Context, blobsList []queue.BlobInfo) ([]queue.BlobInfo, error)
	AssignTasksToWorker(ctx context.Context, workerID string, opt AssignTasksOptions) error
	WorkerTasks(ctx context.Context, workerID string) ([]queue.Job, error)
	MarkTask(ctx context.Context, item queue.Job, status queue.JobStatus) error
}

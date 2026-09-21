package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*CacheCleanupTask)(nil)

const CacheCleanup = "cache-cleanup"

type CacheCleanupTask struct {
	jobService *service.JobService
}

func NewCacheCleanupTask(jobService *service.JobService) *CacheCleanupTask {
	return &CacheCleanupTask{
		jobService: jobService,
	}
}

func (j *CacheCleanupTask) Info() service.TaskInfo {
	return service.TaskInfo{
		Name:        CacheCleanup,
		DisplayName: "Cache Cleanup",
		Schedule:    "",
	}
}

func (j *CacheCleanupTask) Run(ctx context.Context) error {
	return j.jobService.PushJob(
		ctx,
		jobs.CacheCleanup,
		nil,
	)
}

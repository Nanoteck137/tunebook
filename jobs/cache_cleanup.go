package jobs

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*CacheCleanupJob)(nil)

const CacheCleanup = "cache-cleanup"

type CacheCleanupJob struct {
	filesystem *service.FilesystemService
}

func NewCacheCleanupJob(filesystem *service.FilesystemService) *CacheCleanupJob {
	return &CacheCleanupJob{
		filesystem: filesystem,
	}
}

func (j *CacheCleanupJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        CacheCleanup,
		DisplayName: "Cache Cleanup",
	}
}

func (j *CacheCleanupJob) Run(ctx context.Context, data string) error {
	return j.filesystem.ClearCache()
}

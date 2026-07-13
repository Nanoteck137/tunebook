package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*CacheCleanupTask)(nil)

type CacheCleanupTask struct {
	filesystem *service.FilesystemService
}

func NewCacheCleanupTask(filesystem *service.FilesystemService) *CacheCleanupTask {
	return &CacheCleanupTask{
		filesystem: filesystem,
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
	return j.filesystem.ClearCache()
}

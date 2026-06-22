package tasks

import (
	"context"
	"os"

	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/types"
	"github.com/nanoteck137/tunebook/utils"
)

var _ service.Task = (*CacheCleanupTask)(nil)

type CacheCleanupTask struct {
	dataDir types.DataDir
}

func NewCacheCleanupTask(dataDir types.DataDir) *CacheCleanupTask {
	return &CacheCleanupTask{
		dataDir: dataDir,
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
	cacheDir := j.dataDir.Cache()

	err := os.RemoveAll(cacheDir)
	if err != nil {
		return err
	}

	err = utils.CreateDirectories([]string{
		cacheDir,
	})
	if err != nil {
		return err
	}

	return nil
}

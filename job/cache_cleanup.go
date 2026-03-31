package job

import (
	"context"
	"os"

	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/tools/utils"
	"github.com/nanoteck137/tunebook/types"
)

var _ service.Job = (*CacheCleanupJob)(nil)

type CacheCleanupJob struct {
	dataDir types.DataDir
}

func NewCacheCleanupJob(dataDir types.DataDir) *CacheCleanupJob {
	return &CacheCleanupJob{
		dataDir: dataDir,
	}
}

func (j *CacheCleanupJob) Name() string {
	return CacheCleanup
}

func (j *CacheCleanupJob) Schedule() string {
	return ""
}

func (j *CacheCleanupJob) Run(ctx context.Context) error {
	cacheDir := j.dataDir.Cache()

	err := os.RemoveAll(cacheDir.String())
	if err != nil {
		return err
	}

	err = utils.CreateDirectories([]string{
		cacheDir.String(),
	})
	if err != nil {
		return err
	}

	return nil
}

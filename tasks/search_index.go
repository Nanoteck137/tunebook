package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/jobs"
	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*SearchIndexTask)(nil)

const SearchIndex = "search-index"

type SearchIndexTask struct {
	jobService *service.JobService
}

func NewSearchIndexTask(jobService *service.JobService) *SearchIndexTask {
	return &SearchIndexTask{
		jobService: jobService,
	}
}

func (j *SearchIndexTask) Info() service.TaskInfo {
	return service.TaskInfo{
		Name:        SearchIndex,
		DisplayName: "Search Index",
		Schedule:    "",
	}
}

func (j *SearchIndexTask) Run(ctx context.Context) error {
	// Shares the unique key with the post-sync and post-cleanup dispatches so
	// a scheduled run cannot race a reindex that was just queued because the
	// library changed.
	return jobs.DispatchSearchIndex(ctx, j.jobService)
}

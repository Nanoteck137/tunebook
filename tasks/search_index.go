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
	return j.jobService.PushJob(
		ctx,
		jobs.SearchIndex,
		nil,
	)
}

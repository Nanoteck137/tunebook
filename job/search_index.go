package job

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*SearchIndexJob)(nil)

type SearchIndexJob struct {
	searchService *service.SearchService
}

func NewSearchIndexJob(searchService *service.SearchService) *SearchIndexJob {
	return &SearchIndexJob{
		searchService: searchService,
	}
}

func (j *SearchIndexJob) Name() string {
	return SearchIndex
}

func (j *SearchIndexJob) Schedule() string {
	return ""
}

func (j *SearchIndexJob) Run(ctx context.Context) error {
	return j.searchService.Index(ctx)
}

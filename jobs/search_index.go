package jobs

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*SearchIndexJob)(nil)

const SearchIndex = "search-index"

type SearchIndexJob struct {
	searchService *service.SearchService
}

func NewSearchIndexJob(searchService *service.SearchService) *SearchIndexJob {
	return &SearchIndexJob{
		searchService: searchService,
	}
}

func (j *SearchIndexJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        SearchIndex,
		DisplayName: "Search Index",
	}
}

func (j *SearchIndexJob) Run(ctx context.Context, data string) error {
	return j.searchService.Index(ctx)
}

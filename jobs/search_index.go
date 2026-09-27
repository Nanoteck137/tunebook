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
		NoTimeout:   true,
	}
}

func (j *SearchIndexJob) Run(ctx context.Context, data string) error {
	return j.searchService.Index(ctx)
}

// DispatchSearchIndex queues a search reindex. It is used after operations
// that change the library, so the search index does not lag behind the data.
//
// The unique key keeps repeated syncs from piling up duplicate reindex jobs:
// while one is pending or running, further requests are absorbed by it.
func DispatchSearchIndex(
	ctx context.Context,
	jobService *service.JobService,
) error {
	return jobService.PushJob(
		ctx,
		SearchIndex,
		nil,
		service.WithUniqueKey(SearchIndex),
	)
}

package tasks

import (
	"context"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Task = (*SearchIndexTask)(nil)

type SearchIndexTask struct {
	searchService *service.SearchService
}

func NewSearchIndexTask(searchService *service.SearchService) *SearchIndexTask {
	return &SearchIndexTask{
		searchService: searchService,
	}
}

func (j *SearchIndexTask) Name() string {
	return SearchIndex
}

func (j *SearchIndexTask) Schedule() string {
	return ""
}

func (j *SearchIndexTask) Run(ctx context.Context) error {
	return j.searchService.Index(ctx)
}

package jobs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*JobsCleanupJob)(nil)

const JobsCleanup = "jobs-cleanup"

var defaultJobsRetention = 7 * 24 * time.Hour

type JobsCleanupJob struct {
	jobService *service.JobService
}

func NewJobsCleanupJob(jobService *service.JobService) *JobsCleanupJob {
	return &JobsCleanupJob{
		jobService: jobService,
	}
}

func (j *JobsCleanupJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        JobsCleanup,
		DisplayName: "Jobs Cleanup",
	}
}

func (j *JobsCleanupJob) Run(ctx context.Context, data string) error {
	var params service.JobsCleanupParams
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		return err
	}

	retention := time.Duration(params.OlderThanMs) * time.Millisecond
	if retention <= 0 {
		retention = defaultJobsRetention
	}

	_, err = j.jobService.CleanupOldJobs(ctx, retention)
	return err
}

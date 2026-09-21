package jobs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*TestJob)(nil)

const Test = "test"

const defaultTestJobDuration = 30 * time.Second

type TestJob struct{}

func NewTestJob() *TestJob {
	return &TestJob{}
}

func (j *TestJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        Test,
		DisplayName: "Test Job",
	}
}

func (j *TestJob) Run(ctx context.Context, data string) error {
	var params service.TestJobParams
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		return err
	}

	duration := time.Duration(params.DurationMs) * time.Millisecond
	if duration <= 0 {
		duration = defaultTestJobDuration
	}

	select {
	case <-time.After(duration):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

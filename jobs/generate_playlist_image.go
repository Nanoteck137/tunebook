package jobs

import (
	"context"
	"encoding/json"

	"github.com/nanoteck137/tunebook/service"
)

var _ service.Job = (*GeneratePlaylistImageJob)(nil)

const GeneratePlaylistImage = "generate-playlist-image"

type GeneratePlaylistImageJob struct {
	playlistService *service.PlaylistService
}

func NewGeneratePlaylistImageJob(
	playlistService *service.PlaylistService,
) *GeneratePlaylistImageJob {
	return &GeneratePlaylistImageJob{
		playlistService: playlistService,
	}
}

func (j *GeneratePlaylistImageJob) Info() service.JobInfo {
	return service.JobInfo{
		Name:        GeneratePlaylistImage,
		DisplayName: "Generate Playlist Image",
	}
}

func (j *GeneratePlaylistImageJob) Run(ctx context.Context, data string) error {
	var params service.GeneratePlaylistImageParams
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		return err
	}

	return j.playlistService.GeneratePlaylistImage(ctx, params)
}

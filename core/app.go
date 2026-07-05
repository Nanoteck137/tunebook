package core

import (
	"github.com/nanoteck137/tunebook/config"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/service"
	"github.com/nanoteck137/tunebook/tools/broker"
	"github.com/nanoteck137/tunebook/types"
)

// Inspiration from Pocketbase: https://github.com/pocketbase/pocketbase
// File: https://github.com/pocketbase/pocketbase/blob/master/core/app.go
type App interface {
	DB() *database.Database
	Config() *config.Config

	DataDir() types.DataDir

	NotificationService() *service.NotificationService
	TaskService() *service.TaskService
	JobQueueService() *service.JobQueueService
	AuthService() *service.AuthService
	UserService() *service.UserService
	SearchService() *service.SearchService
	LibraryService() *service.LibraryService
	ImageService() *service.ImageService
	MediaService() *service.MediaService

	ArtistService() *service.ArtistService
	AlbumService() *service.AlbumService
	TrackService() *service.TrackService
	PlaylistService() *service.PlaylistService
	HistoryService() *service.HistoryService
	QueueService() *service.QueueService

	Broker() *broker.Broker

	Bootstrap() error
	Shutdown() error
}

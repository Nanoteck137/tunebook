package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gosimple/slug"
	"github.com/nanoteck137/tunebook"
	"github.com/nanoteck137/tunebook/config"
)

var notificationErr = NewServiceErrCreator("notification")

type NotifyPriority int

const (
	NotifyPriorityMin    NotifyPriority = 1
	NotifyPriorityLow    NotifyPriority = 2
	NotifyPriorityNormal NotifyPriority = 3
	NotifyPriorityHigh   NotifyPriority = 4
	NotifyPriorityUrgent NotifyPriority = 5
)

func (p NotifyPriority) String() string {
	switch p {
	case NotifyPriorityMin:
		return "min"
	case NotifyPriorityLow:
		return "low"
	case NotifyPriorityNormal:
		return "normal"
	case NotifyPriorityHigh:
		return "high"
	case NotifyPriorityUrgent:
		return "urgent"
	default:
		return "normal"
	}
}

const (
	NotifyTagWarning       = "warning"
	NotifyTagRotatingLight = "rotating_light"
	NotifyTagComputer      = "computer"
	NotifyTagCd            = "cd"
	NotifyTagLoadspeaker   = "loudspeaker"
)

type NotificationService struct {
	logger *slog.Logger

	BaseUrl string
	Topic   string
}

func NewNotificationService(
	logger *slog.Logger, 
	config *config.Config,
) *NotificationService {
	return &NotificationService{
		logger:  logger,
		BaseUrl: config.NtfyBaseUrl,
		Topic:   config.NtfyTopic,
	}
}

func (s *NotificationService) IsEnabled() bool {
	return s.BaseUrl != "" && s.Topic != ""
}

type notificationMessageBody struct {
	Topic    string   `json:"topic"`
	Message  string   `json:"message,omitempty"`
	Title    string   `json:"title"`
	Tags     []string `json:"tags,omitempty"`
	Priority int      `json:"priority,omitempty"`

	Markdown bool `json:"markdown"`
}

func (s *NotificationService) sendNotification(
	body notificationMessageBody,
) error {
	d, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequest("POST", s.BaseUrl, bytes.NewReader(d))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("returned non 200 status code: %d", resp.StatusCode)
	}

	return nil
}

type SendSimpleParams struct {
	Title    string
	Message  string
	Tags     []string
	Priority NotifyPriority
}

func (s *NotificationService) SendSimple(params SendSimpleParams) error {
	if !s.IsEnabled() {
		s.logger.Debug("notification service not configured, skipping")
		return nil
	}

	if params.Title == "" {
		return notificationErr.New("send simple: no title set")
	}

	if params.Priority == 0 {
		params.Priority = NotifyPriorityNormal
	}

	params.Tags = append(params.Tags, slug.Make(tunebook.AppName))

	body := notificationMessageBody{
		Topic:    s.Topic,
		Message:  params.Message,
		Title:    tunebook.AppName + ": " + params.Title,
		Tags:     params.Tags,
		Priority: int(params.Priority),
		Markdown: true,
	}

	err := s.sendNotification(body)
	if err != nil {
		s.logger.Error("simple notification",
			slog.String("title", params.Title),
			slog.String("message", params.Message),
			slog.String("priority", params.Priority.String()),
			slog.Any("tags", params.Tags),
			slog.Any("err", err),
		)

		return notificationErr.Wrap("send simple", err)
	}

	s.logger.Debug("send simple notification",
		slog.String("title", params.Title),
		slog.String("message", params.Message),
		slog.String("priority", params.Priority.String()),
		slog.Any("tags", params.Tags),
	)

	return nil
}

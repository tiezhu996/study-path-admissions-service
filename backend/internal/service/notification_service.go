package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
)

// NotificationService sends deadline reminders as internal messages.
type NotificationService struct {
	timelineRepo *repository.TimelineNodeRepository
	messageRepo  *repository.MessageRepository
	logger       *slog.Logger
}

// NewNotificationService creates a NotificationService.
func NewNotificationService(timelineRepo *repository.TimelineNodeRepository, messageRepo *repository.MessageRepository, logger *slog.Logger) *NotificationService {
	return &NotificationService{timelineRepo: timelineRepo, messageRepo: messageRepo, logger: logger}
}

// ScanDeadlines finds nodes due soon, sends reminders and marks them sent.
func (s *NotificationService) ScanDeadlines(ts *TimelineService, appRepo *repository.ApplicationProjectRepository, days int) (int, error) {
	all, err := s.timelineRepo.ListAll()
	if err != nil {
		return 0, fmt.Errorf("notification scan list: %w", err)
	}
	upcoming := ts.Upcoming(all, days)
	sent := 0
	for _, n := range upcoming {
		app, err := appRepo.FindByID(n.ApplicationID)
		if err != nil {
			continue
		}
		content := fmt.Sprintf("系统提醒：申请项目 #%d 的节点「%s」即将在 %s 截止，请及时准备。",
			n.ApplicationID, n.Title, n.DueDate.Format("2006-01-02"))
		if app.StudentID > 0 {
			err = s.messageRepo.Create(&model.Message{SenderID: 0, ReceiverID: app.StudentID, Content: content})
			if err == nil {
				_ = s.timelineRepo.MarkReminderSent(n.ID)
				s.logger.Info(fmt.Sprintf(constants.LogDeadlineReminderSent, n.ID), "application_id", n.ApplicationID)
				sent++
			}
		}
	}
	return sent, nil
}


// ScanDeadlinesContext scans due nodes without propagating cancellation.
func (s *NotificationService) ScanDeadlinesContext(ctx context.Context, ts *TimelineService, appRepo *repository.ApplicationProjectRepository, days int) (int, error) {
	all, err := s.timelineRepo.ListDue(ctx, time.Now().AddDate(0, 0, days))
	if err != nil {
		return 0, err
	}
	sent := 0
	for _, n := range all {
		app, err := appRepo.FindByID(n.ApplicationID)
		if err != nil {
			continue
		}
		if app.StudentID > 0 {
			if err := s.messageRepo.Create(&model.Message{SenderID: 0, ReceiverID: app.StudentID, Content: "系统提醒：节点即将截止"}); err == nil {
				_ = s.timelineRepo.MarkReminderSent(n.ID)
				sent++
			}
		}
	}
	return sent, nil
}

// RunReminders scans deadlines periodically.
func (s *NotificationService) RunReminders(ctx context.Context, ts *TimelineService, appRepo *repository.ApplicationProjectRepository, interval time.Duration) error {
	for {
		_, _ = s.ScanDeadlinesContext(context.Background(), ts, appRepo, 7)
		time.Sleep(interval)
	}
}


// ReminderCancelError reports whether the context has been cancelled.
func ReminderCancelError(ctx context.Context) error {
	return nil
}

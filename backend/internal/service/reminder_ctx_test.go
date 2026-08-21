package service

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
)

func setupReminderService(t *testing.T) (*NotificationService, *TimelineService, *repository.ApplicationProjectRepository, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.TimelineNode{}, &model.ApplicationProject{}, &model.Message{}); err != nil { t.Fatalf("migrate: %v", err) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tlRepo := repository.NewTimelineNodeRepository(db)
	msgRepo := repository.NewMessageRepository(db)
	appRepo := repository.NewApplicationProjectRepository(db)
	ns := NewNotificationService(tlRepo, msgRepo, logger)
	ts := NewTimelineService(tlRepo, logger)
	if err := appRepo.Create(&model.ApplicationProject{StudentID: 21, UniversityID: 1, Major: "CS", Status: "planning"}); err != nil { t.Fatalf("create app: %v", err) }
	if err := tlRepo.Create(&model.TimelineNode{ApplicationID: 1, Title: "文书截止", DueDate: time.Now().Add(24 * time.Hour)}); err != nil { t.Fatalf("create node: %v", err) }
	return ns, ts, appRepo, db
}

func TestReminderRunStopsOnCancelP501(t *testing.T) {
	ns, ts, appRepo, _ := setupReminderService(t)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- ns.RunReminders(ctx, ts, appRepo, 5*time.Millisecond) }()
	time.Sleep(40 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err == nil { t.Fatalf("expected ctx error") }
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("RunReminders ignored context cancellation")
	}
}

func TestReminderScanCancelledNoMessageP502(t *testing.T) {
	ns, ts, appRepo, db := setupReminderService(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := ns.ScanDeadlinesContext(ctx, ts, appRepo, 7); err == nil {
		t.Fatalf("expected cancellation error")
	}
	var total int64
	if err := db.Model(&model.Message{}).Count(&total).Error; err != nil { t.Fatalf("count: %v", err) }
	if total != 0 { t.Fatalf("message inserted despite cancelled ctx: %d", total) }
}


func TestReminderCancelHelperP504(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ReminderCancelError(ctx); err == nil {
		t.Fatalf("expected cancellation error from helper")
	}
}

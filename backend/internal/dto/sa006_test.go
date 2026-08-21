package dto_test

import (
	"io"
	"log/slog"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
)

func TestSA006MessageBatch(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memdb?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(4) }
	if err := db.AutoMigrate(&model.Message{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := repository.NewMessageRepository(db)
	svc := service.NewMessageService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))
	count, err := svc.SendBatch([]uint{11, 12, 13}, "材料已更新")
	if err != nil { t.Fatalf("send batch: %v", err) }
	if count != 3 { t.Fatalf("sent count = %d, want 3", count) }
	for _, id := range []uint{11, 12, 13} {
		items, err := svc.ListByReceiver(id, false)
		if err != nil { t.Fatalf("list receiver %d: %v", id, err) }
		if len(items) != 1 { t.Fatalf("receiver %d has %d messages, want 1", id, len(items)) }
	}
}

package repository

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

func TestTimelineListDueCancelledP503(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.TimelineNode{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := NewTimelineNodeRepository(db)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := repo.ListDue(ctx, time.Now()); err == nil {
		t.Fatalf("expected cancellation error")
	}
}

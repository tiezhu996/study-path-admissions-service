package repository

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

func TestRecommendationRepoPersistsStatusP704(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.Recommendation{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := NewRecommendationRepository(db)
	if err := repo.Create(&model.Recommendation{StudentID: 9, CounselorID: 2, UniversityIDs: "[1]", Reason: "冲刺", Status: "sent"}); err != nil { t.Fatalf("create: %v", err) }
	rec, err := repo.FindByID(1)
	if err != nil { t.Fatalf("find: %v", err) }
	rec.Status = "accepted"
	if err := repo.UpdateStatus(rec); err != nil { t.Fatalf("update status: %v", err) }
	got, err := repo.FindByID(rec.ID)
	if err != nil { t.Fatalf("find: %v", err) }
	if got.Status != "accepted" {
		t.Fatalf("persisted status = %q, want accepted", got.Status)
	}
}

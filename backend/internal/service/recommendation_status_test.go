package service

import (
	"io"
	"log/slog"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
)

func setupRecommendationService(t *testing.T) (*RecommendationService, *repository.RecommendationRepository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.Recommendation{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := repository.NewRecommendationRepository(db)
	svc := NewRecommendationService(repo, repository.NewUniversityRepository(db), slog.New(slog.NewTextHandler(io.Discard, nil)))
	return svc, repo
}

func TestRecommendationTransitionSentAcceptedP701(t *testing.T) {
	svc, repo := setupRecommendationService(t)
	rec, err := svc.Create(2, 9, []uint{1}, "冲刺校")
	if err != nil { t.Fatalf("create: %v", err) }
	rec.Status = "sent"
	if err := repo.UpdateStatus(rec); err != nil { t.Fatalf("set sent: %v", err) }
	got, err := svc.UpdateStatus(rec.ID, "accepted")
	if err != nil { t.Fatalf("sent->accepted: %v", err) }
	if got.Status != "accepted" { t.Fatalf("status = %q, want accepted", got.Status) }
}

func TestRecommendationIllegalTransitionP702(t *testing.T) {
	svc, _ := setupRecommendationService(t)
	rec, err := svc.Create(2, 9, []uint{1}, "保底校")
	if err != nil { t.Fatalf("create: %v", err) }
	if _, err := svc.UpdateStatus(rec.ID, "accepted"); err == nil {
		t.Fatalf("draft->accepted should be rejected")
	}
}

package model_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
)

func sa007Setup(t *testing.T) (*service.RecommendationService, *repository.RecommendationRepository, *handler.RecommendationHandler) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.Recommendation{}); err != nil { t.Fatalf("migrate: %v", err) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := repository.NewRecommendationRepository(db)
	svc := service.NewRecommendationService(repo, repository.NewUniversityRepository(db), logger)
	return svc, repo, handler.NewRecommendationHandler(svc, logger)
}

func TestSA007RecTransition(t *testing.T) {
	svc, repo, _ := sa007Setup(t)
	rec, err := svc.Create(2, 9, []uint{1}, "冲刺校")
	if err != nil { t.Fatalf("create: %v", err) }
	rec.Status = "sent"
	if err := repo.UpdateStatus(rec); err != nil { t.Fatalf("set sent: %v", err) }
	got, err := svc.UpdateStatus(rec.ID, "accepted")
	if err != nil { t.Fatalf("sent->accepted: %v", err) }
	if got.Status != "accepted" { t.Fatalf("status = %q, want accepted", got.Status) }
}

func TestSA007RecHandler(t *testing.T) {
	_, repo, h := sa007Setup(t)
	if err := repo.Create(&model.Recommendation{StudentID: 9, CounselorID: 2, UniversityIDs: "[1]", Reason: "保底", Status: "draft"}); err != nil { t.Fatalf("create: %v", err) }
	rec, err := repo.FindByID(1)
	if err != nil { t.Fatalf("find: %v", err) }
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.PUT("/recommendations/:id/status", h.UpdateStatus)
	req := httptest.NewRequest(http.MethodPut, "/recommendations/"+strconv.FormatUint(uint64(rec.ID), 10)+"/status", strings.NewReader(`{"status":"accepted"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusOK { t.Fatalf("illegal transition returned 200; body=%s", w.Body.String()) }
}

func TestSA007RecRepo(t *testing.T) {
	_, repo, _ := sa007Setup(t)
	if err := repo.Create(&model.Recommendation{StudentID: 9, CounselorID: 2, UniversityIDs: "[1]", Reason: "冲刺", Status: "sent"}); err != nil { t.Fatalf("create: %v", err) }
	rec, err := repo.FindByID(1)
	if err != nil { t.Fatalf("find: %v", err) }
	rec.Status = "accepted"
	if err := repo.UpdateStatus(rec); err != nil { t.Fatalf("update status: %v", err) }
	got, err := repo.FindByID(rec.ID)
	if err != nil { t.Fatalf("find again: %v", err) }
	if got.Status != "accepted" { t.Fatalf("persisted status = %q, want accepted", got.Status) }
}

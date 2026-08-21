package handler

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

	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
)

func setupRecHandler(t *testing.T) (*RecommendationHandler, *repository.RecommendationRepository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.Recommendation{}); err != nil { t.Fatalf("migrate: %v", err) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := repository.NewRecommendationRepository(db)
	svc := service.NewRecommendationService(repo, repository.NewUniversityRepository(db), logger)
	return NewRecommendationHandler(svc, logger), repo
}

func TestRecommendationHandlerIllegalStatusP703(t *testing.T) {
	h, repo := setupRecHandler(t)
	if err := repo.Create(&model.Recommendation{StudentID: 9, CounselorID: 2, UniversityIDs: "[1]", Reason: "保底", Status: "draft"}); err != nil { t.Fatalf("create: %v", err) }
	rec, err := repo.FindByID(1)
	if err != nil { t.Fatalf("find: %v", err) }
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.PUT("/recommendations/:id/status", h.UpdateStatus)
	req := httptest.NewRequest(http.MethodPut, "/recommendations/"+itoa(rec.ID)+"/status", strings.NewReader(`{"status":"accepted"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		t.Fatalf("illegal transition returned 200; body=%s", w.Body.String())
	}
}

func itoa(v uint) string { return strconv.FormatUint(uint64(v), 10) }

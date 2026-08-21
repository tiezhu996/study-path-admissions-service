package handler

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
)

func TestUniversityHandlerSnapshotStableP405(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.University{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := repository.NewUniversityRepository(db)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	if err := repo.Create(&model.University{Name:"Harvard", Country:"US", City:"Boston", Ranking:1}); err != nil { t.Fatalf("create us: %v", err) }
	if err := repo.Create(&model.University{Name:"Oxford", Country:"UK", City:"Oxford", Ranking:2}); err != nil { t.Fatalf("create uk: %v", err) }
	h := NewUniversityHandler(service.NewUniversityService(repo, logger), logger)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/universities", h.List)
	first := httptest.NewRecorder()
	r.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/universities?country=US", nil))
	second := httptest.NewRecorder()
	r.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/universities?country=UK", nil))
	if first.Body.String() != first.Body.String() { t.Fatalf("nop") }
	// First response must still contain Harvard after the second call.
	if !strings.Contains(first.Body.String(), "Harvard") {
		t.Fatalf("first response lost Harvard after second query: %s", first.Body.String())
	}
	if !strings.Contains(second.Body.String(), "Oxford") {
		t.Fatalf("second response missing Oxford: %s", second.Body.String())
	}
}

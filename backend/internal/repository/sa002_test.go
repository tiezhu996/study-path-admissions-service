package repository_test

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

func sa002Svc(t *testing.T) *service.DocumentService {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.Document{}, &model.DocumentVersion{}, &model.Annotation{}, &model.ApplicationProject{}); err != nil { t.Fatalf("migrate: %v", err) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return service.NewDocumentService(db, repository.NewDocumentRepository(db), repository.NewDocumentVersionRepository(db), repository.NewAnnotationRepository(db), repository.NewApplicationProjectRepository(db), logger)
}

func TestSA002DocRollback404(t *testing.T) {
	svc := sa002Svc(t)
	h := handler.NewDocumentHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.POST("/documents/:id/rollback", h.Rollback)
	req := httptest.NewRequest(http.MethodPost, "/documents/8/rollback", strings.NewReader(`{"version_no":3}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound { t.Fatalf("status = %d, want 404; body=%s", w.Code, w.Body.String()) }
}

func TestSA002DocGet404(t *testing.T) {
	svc := sa002Svc(t)
	h := handler.NewDocumentHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.GET("/documents/:id", h.Get)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/documents/8", nil))
	if w.Code != http.StatusNotFound { t.Fatalf("status = %d, want 404; body=%s", w.Code, w.Body.String()) }
}

func TestSA002VersionChain(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.DocumentVersion{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := repository.NewDocumentVersionRepository(db)
	if _, err := repo.FindByDocumentAndVersion(99, 3); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, got %v", err)
	}
}

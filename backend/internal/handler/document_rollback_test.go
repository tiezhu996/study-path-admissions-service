package handler

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

	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
)

func setupDocService(t *testing.T) (*service.DocumentService, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if err := db.AutoMigrate(&model.Document{}, &model.DocumentVersion{}, &model.Annotation{}, &model.ApplicationProject{}); err != nil { t.Fatalf("migrate: %v", err) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	docRepo := repository.NewDocumentRepository(db)
	verRepo := repository.NewDocumentVersionRepository(db)
	annRepo := repository.NewAnnotationRepository(db)
	appRepo := repository.NewApplicationProjectRepository(db)
	svc := service.NewDocumentService(db, docRepo, verRepo, annRepo, appRepo, logger)
	return svc, db
}

func TestDocumentRollbackMissing404P201(t *testing.T) {
	svc, _ := setupDocService(t)
	h := NewDocumentHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.POST("/documents/:id/rollback", h.Rollback)
	req := httptest.NewRequest(http.MethodPost, "/documents/8/rollback", strings.NewReader(`{"version_no":3}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusNotFound, w.Body.String())
	}
	_ = errors.Is
	_ = constants.CodeNotFound
}


func TestDocumentGetMissing404P203(t *testing.T) {
	svc, _ := setupDocService(t)
	h := NewDocumentHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.GET("/documents/:id", h.Get)
	req := httptest.NewRequest(http.MethodGet, "/documents/8", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusNotFound, w.Body.String())
	}
}

package handler_test

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

	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
)

func sa003Setup(t *testing.T) (*service.MaterialService, *repository.MaterialItemRepository, *handler.MaterialHandler) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.MaterialItem{}); err != nil { t.Fatalf("migrate: %v", err) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := repository.NewMaterialItemRepository(db)
	svc := service.NewMaterialService(repo, logger)
	return svc, repo, handler.NewMaterialHandler(svc, logger)
}

func TestSA003MatProgress(t *testing.T) {
	svc, repo, _ := sa003Setup(t)
	if err := repo.Create(&model.MaterialItem{ApplicationID: 5, Name: "成绩单", Status: "uploaded"}); err != nil { t.Fatalf("create item: %v", err) }
	if _, err := svc.Progress(5); err != nil { t.Fatalf("progress: %v", err) }
}

func TestSA003MatStatusCache(t *testing.T) {
	svc, repo, _ := sa003Setup(t)
	if err := repo.Create(&model.MaterialItem{ApplicationID: 5, Name: "护照", Status: "pending"}); err != nil { t.Fatalf("create item: %v", err) }
	item, err := repo.FindByID(1)
	if err != nil { t.Fatalf("find item: %v", err) }
	if _, err := svc.UpdateStatus(9, item.ID, "student", "uploaded", "http://file/1"); err != nil { t.Fatalf("update status: %v", err) }
}

func TestSA003MatHandlerCreate(t *testing.T) {
	_, _, h := sa003Setup(t)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.POST("/applications/:id/materials", h.Create)
	req := httptest.NewRequest(http.MethodPost, "/applications/5/materials", strings.NewReader(`{"name":"成绩单","category":"成绩"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated { t.Fatalf("status = %d, body=%s", w.Code, w.Body.String()) }
}

func TestSA003MatTouchIndex(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.MaterialItem{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := repository.NewMaterialItemRepository(db)
	if err := repo.TouchProgressIndex(7); err != nil { t.Fatalf("touch: %v", err) }
}

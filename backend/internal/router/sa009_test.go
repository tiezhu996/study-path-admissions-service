package router_test

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

func sa009User(t *testing.T) (*service.UserService, *repository.UserRepository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.User{}); err != nil { t.Fatalf("migrate: %v", err) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := repository.NewUserRepository(db)
	cfg := &config.Config{JWTSecret:"test-secret", JWTExpire: time.Hour}
	return service.NewUserService(repo, logger, cfg), repo
}

func TestSA009Login401(t *testing.T) {
	svc, _ := sa009User(t)
	h := handler.NewUserHandler(svc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.POST("/users/login", h.Login)
	req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"username":"ghost","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized { t.Fatalf("status = %d, want 401; body=%s", w.Code, w.Body.String()) }
}

func TestSA009UserChain(t *testing.T) {
	_, repo := sa009User(t)
	if _, err := repo.FindByUsername("ghost"); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, got %v", err)
	}
}

func TestSA009UserGet404(t *testing.T) {
	svc, _ := sa009User(t)
	_, err := svc.GetByID(404)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected AppError 404, got %v", err)
	}
}

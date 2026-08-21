package service

import (
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

func TestUserGetByIDMissing404P903(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.User{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := repository.NewUserRepository(db)
	cfg := &config.Config{JWTSecret:"test", JWTExpire: time.Hour}
	svc := NewUserService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)), cfg)
	_, err = svc.GetByID(404)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected AppError 404, got %v", err)
	}
}

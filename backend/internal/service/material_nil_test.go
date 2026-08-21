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

func setupMaterialService(t *testing.T) (*MaterialService, *repository.MaterialItemRepository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.MaterialItem{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := repository.NewMaterialItemRepository(db)
	svc := NewMaterialService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))
	return svc, repo
}

func TestMaterialProgressNoPanicP301(t *testing.T) {
	svc, repo := setupMaterialService(t)
	if err := repo.Create(&model.MaterialItem{ApplicationID: 5, Name: "成绩单", Status: "uploaded"}); err != nil { t.Fatalf("create item: %v", err) }
	if _, err := svc.Progress(5); err != nil { t.Fatalf("progress: %v", err) }
}

func TestMaterialStatusCacheNoPanicP306(t *testing.T) {
	svc, repo := setupMaterialService(t)
	if err := repo.Create(&model.MaterialItem{ApplicationID: 5, Name: "护照", Status: "pending"}); err != nil { t.Fatalf("create item: %v", err) }
	item, err := repo.FindByID(1)
	if err != nil { t.Fatalf("find item: %v", err) }
	if _, err := svc.UpdateStatus(9, item.ID, "student", "uploaded", "http://file/1"); err != nil { t.Fatalf("update status: %v", err) }
}

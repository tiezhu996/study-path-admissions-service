package repository

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

func TestMaterialTouchProgressIndexNoPanicP305(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.MaterialItem{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := NewMaterialItemRepository(db)
	if err := repo.TouchProgressIndex(7); err != nil { t.Fatalf("touch: %v", err) }
}

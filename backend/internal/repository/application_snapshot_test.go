package repository

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

func TestApplicationRepoListAllSnapshotStableP103(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.ApplicationProject{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := NewApplicationProjectRepository(db)
	if err := repo.Create(&model.ApplicationProject{StudentID: 1, UniversityID: 1, Major: "A", Status: "planning"}); err != nil { t.Fatalf("create a: %v", err) }
	if err := repo.Create(&model.ApplicationProject{StudentID: 2, UniversityID: 1, Major: "B", Status: "planning"}); err != nil { t.Fatalf("create b: %v", err) }
	first, err := repo.ListAll()
	if err != nil { t.Fatalf("first list: %v", err) }
	if len(first) != 2 { t.Fatalf("first len = %d", len(first)) }
	second, err := repo.ListAll()
	if err != nil { t.Fatalf("second list: %v", err) }
	second[0].Major = "CHANGED"
	if first[0].Major == "CHANGED" { t.Fatalf("first snapshot mutated") }
}

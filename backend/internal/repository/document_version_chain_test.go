package repository

import (
	"errors"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

func TestDocumentVersionRepoChainErrNotFoundP202(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if err := db.AutoMigrate(&model.DocumentVersion{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := NewDocumentVersionRepository(db)
	_, err = repo.FindByDocumentAndVersion(99, 3)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, got %v", err)
	}
}

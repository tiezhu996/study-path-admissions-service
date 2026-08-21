package repository

import (
	"errors"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

func TestUserRepoChainErrNotFoundP902(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if err := db.AutoMigrate(&model.User{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := NewUserRepository(db)
	_, err = repo.FindByUsername("ghost")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, got %v", err)
	}
}

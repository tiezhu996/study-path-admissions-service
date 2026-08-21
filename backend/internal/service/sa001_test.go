package service_test

import (
	"io"
	"log/slog"
	"sync"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
)

func sa001DB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open sqlite: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.University{}, &model.ApplicationProject{}); err != nil { t.Fatalf("migrate: %v", err) }
	return db
}

func TestSA001AppCacheRace(t *testing.T) {
	db := sa001DB(t)
	univRepo := repository.NewUniversityRepository(db)
	if err := univRepo.Create(&model.University{Name: "State University", Country: "US"}); err != nil { t.Fatalf("create university: %v", err) }
	repo := repository.NewApplicationProjectRepository(db)
	for i := 0; i < 12; i++ {
		if err := repo.Create(&model.ApplicationProject{StudentID: 7, UniversityID: 1, Major: "CS", Status: "planning"}); err != nil { t.Fatalf("create project: %v", err) }
	}
	svc := service.NewApplicationService(repo, univRepo, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if _, err := svc.List(7, "student"); err != nil { t.Fatalf("warm cache: %v", err) }
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(1)
	go func() { defer wg.Done(); <-start; for i := 0; i < 200; i++ { _, _ = svc.Get(1, 7, "student") } }()
	wg.Add(1)
	go func() { defer wg.Done(); <-start; for i := 0; i < 200; i++ { _, _ = svc.UpdateStatus(uint((i%12)+1), 7, "student", "preparing") } }()
	close(start)
	wg.Wait()
}

func TestSA001AppListSnapshot(t *testing.T) {
	db := sa001DB(t)
	univRepo := repository.NewUniversityRepository(db)
	repo := repository.NewApplicationProjectRepository(db)
	if err := univRepo.Create(&model.University{Name: "Tech University", Country: "US"}); err != nil { t.Fatalf("create university: %v", err) }
	if err := repo.Create(&model.ApplicationProject{StudentID: 8, UniversityID: 1, Major: "Economics", Status: "planning"}); err != nil { t.Fatalf("create project: %v", err) }
	if err := repo.Create(&model.ApplicationProject{StudentID: 8, UniversityID: 1, Major: "Mathematics", Status: "planning"}); err != nil { t.Fatalf("create project: %v", err) }
	svc := service.NewApplicationService(repo, univRepo, slog.New(slog.NewTextHandler(io.Discard, nil)))
	first, err := svc.List(8, "admin")
	if err != nil { t.Fatalf("first list: %v", err) }
	if len(first) != 2 { t.Fatalf("first list len = %d, want 2", len(first)) }
	second, err := svc.List(8, "admin")
	if err != nil { t.Fatalf("second list: %v", err) }
	second[0].Status = "admitted"
	if first[0].Status == "admitted" { t.Fatalf("first snapshot was mutated by second list") }
}

func TestSA001RepoListAllSnapshot(t *testing.T) {
	db := sa001DB(t)
	repo := repository.NewApplicationProjectRepository(db)
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

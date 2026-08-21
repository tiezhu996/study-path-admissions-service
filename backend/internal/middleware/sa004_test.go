package middleware_test

import (
	"io"
	"log/slog"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
)

func sa004Univ(t *testing.T) (*service.UniversityService, *repository.UniversityRepository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.University{}); err != nil { t.Fatalf("migrate: %v", err) }
	repo := repository.NewUniversityRepository(db)
	svc := service.NewUniversityService(repo, slog.New(slog.NewTextHandler(io.Discard, nil)))
	return svc, repo
}

func TestSA004UnivListStable(t *testing.T) {
	svc, repo := sa004Univ(t)
	if err := repo.Create(&model.University{Name:"Harvard", Country:"US", City:"Boston", Ranking:1}); err != nil { t.Fatalf("create us: %v", err) }
	if err := repo.Create(&model.University{Name:"Oxford", Country:"UK", City:"Oxford", Ranking:2}); err != nil { t.Fatalf("create uk: %v", err) }
	first, _, err := svc.List("US", 0, 0, "", 1, 10)
	if err != nil { t.Fatalf("first list: %v", err) }
	if len(first) != 1 || first[0].Name != "Harvard" { t.Fatalf("first = %+v", first) }
	second, _, err := svc.List("UK", 0, 0, "", 1, 10)
	if err != nil { t.Fatalf("second list: %v", err) }
	if len(second) != 1 || second[0].Name != "Oxford" { t.Fatalf("second = %+v", second) }
	if first[0].Name != "Harvard" || first[0].Country != "US" { t.Fatalf("first snapshot mutated: %+v", first[0]) }
}

func TestSA004UnivGetStable(t *testing.T) {
	svc, repo := sa004Univ(t)
	if err := repo.Create(&model.University{Name:"Harvard", Country:"US", City:"Boston", Ranking:1}); err != nil { t.Fatalf("create us: %v", err) }
	if err := repo.Create(&model.University{Name:"Oxford", Country:"UK", City:"Oxford", Ranking:2}); err != nil { t.Fatalf("create uk: %v", err) }
	first, err := svc.Get(1)
	if err != nil { t.Fatalf("first get: %v", err) }
	second, err := svc.Get(2)
	if err != nil { t.Fatalf("second get: %v", err) }
	if first.Name == second.Name { t.Fatalf("first result aliased by second get") }
}

func TestSA004UnivIdsStable(t *testing.T) {
	_, repo := sa004Univ(t)
	if err := repo.Create(&model.University{Name:"MIT", Country:"US", City:"Cambridge", Ranking:1}); err != nil { t.Fatalf("create mit: %v", err) }
	if err := repo.Create(&model.University{Name:"Cambridge", Country:"UK", City:"Cambridge", Ranking:2}); err != nil { t.Fatalf("create cam: %v", err) }
	first, err := repo.ListByIDs([]uint{1})
	if err != nil { t.Fatalf("first ids: %v", err) }
	if len(first) != 1 || first[0].Name != "MIT" { t.Fatalf("first = %+v", first) }
	second, err := repo.ListByIDs([]uint{2})
	if err != nil { t.Fatalf("second ids: %v", err) }
	if len(second) != 1 || second[0].Name != "Cambridge" { t.Fatalf("second = %+v", second) }
	if first[0].Name != "MIT" { t.Fatalf("first snapshot mutated: %+v", first[0]) }
}

func TestSA004UnivRepoListStable(t *testing.T) {
	_, repo := sa004Univ(t)
	if err := repo.Create(&model.University{Name:"Harvard", Country:"US", City:"Boston", Ranking:1}); err != nil { t.Fatalf("create us: %v", err) }
	if err := repo.Create(&model.University{Name:"Oxford", Country:"UK", City:"Oxford", Ranking:2}); err != nil { t.Fatalf("create uk: %v", err) }
	first, _, err := repo.List("US", 0, 0, "", 1, 10)
	if err != nil { t.Fatalf("first list: %v", err) }
	if len(first) != 1 || first[0].Name != "Harvard" { t.Fatalf("first = %+v", first) }
	second, _, err := repo.List("UK", 0, 0, "", 1, 10)
	if err != nil { t.Fatalf("second list: %v", err) }
	if len(second) != 1 || second[0].Name != "Oxford" { t.Fatalf("second = %+v", second) }
	if first[0].Name != "Harvard" { t.Fatalf("first snapshot mutated: %+v", first[0]) }
}

package constants_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
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

func sa008Handler(t *testing.T) *handler.DashboardHandler {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.ApplicationProject{}, &model.MaterialItem{}, &model.University{}); err != nil { t.Fatalf("migrate: %v", err) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	appRepo := repository.NewApplicationProjectRepository(db)
	univRepo := repository.NewUniversityRepository(db)
	matRepo := repository.NewMaterialItemRepository(db)
	if err := univRepo.Create(&model.University{Name:"State", Country:"US"}); err != nil { t.Fatalf("univ: %v", err) }
	for i := 0; i < 8; i++ {
		if err := appRepo.Create(&model.ApplicationProject{StudentID: uint(i+1), UniversityID: 1, Major:"CS", Status:"planning"}); err != nil { t.Fatalf("app: %v", err) }
	}
	return handler.NewDashboardHandler(service.NewApplicationService(appRepo, univRepo, logger), service.NewUniversityService(univRepo, logger), service.NewMaterialService(matRepo, logger), logger)
}

func TestSA008DashConcurrent(t *testing.T) {
	h := sa008Handler(t)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.GET("/stats", h.Stats)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/stats", nil)) }()
	}
	close(start)
	wg.Wait()
}

func TestSA008FillStale(t *testing.T) {
	var dst service.AppStats
	service.FillAppStats([]model.ApplicationProject{{Status: "admitted"}}, &dst)
	service.FillAppStats([]model.ApplicationProject{{Status: "planning"}}, &dst)
	if _, ok := dst.ByStatus["admitted"]; ok { t.Fatalf("stale status leaked: %+v", dst.ByStatus) }
}

func TestSA008FillReset(t *testing.T) {
	var dst service.AppStats
	service.FillAppStats([]model.ApplicationProject{{Status: "admitted"}}, &dst)
	service.FillAppStats([]model.ApplicationProject{{Status: "planning"}}, &dst)
	if dst.Admitted != 0 || dst.Applied != 0 { t.Fatalf("counters not reset: %+v", dst) }
}

func TestSA008ComputeConcurrent(t *testing.T) {
	projects := make([]model.ApplicationProject, 20)
	for i := range projects { projects[i] = model.ApplicationProject{Status: "planning"} }
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; for j := 0; j < 50; j++ { _ = service.ComputeAppStats(projects) } }()
	}
	close(start)
	wg.Wait()
}

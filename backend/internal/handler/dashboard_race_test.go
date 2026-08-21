package handler

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

	"github.com/gbstudyapply/gbstudyapply/internal/middleware"
	"github.com/gbstudyapply/gbstudyapply/internal/model"
	"github.com/gbstudyapply/gbstudyapply/internal/repository"
	"github.com/gbstudyapply/gbstudyapply/internal/service"
)

func setupDashboardHandler(t *testing.T) *DashboardHandler {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil { t.Fatalf("open: %v", err) }
	if sqlDB, err := db.DB(); err == nil { sqlDB.SetMaxOpenConns(1) }
	if err := db.AutoMigrate(&model.ApplicationProject{}, &model.MaterialItem{}, &model.University{}); err != nil { t.Fatalf("migrate: %v", err) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	appRepo := repository.NewApplicationProjectRepository(db)
	univRepo := repository.NewUniversityRepository(db)
	matRepo := repository.NewMaterialItemRepository(db)
	appSvc := service.NewApplicationService(appRepo, univRepo, logger)
	matSvc := service.NewMaterialService(matRepo, logger)
	if err := univRepo.Create(&model.University{Name:"State", Country:"US"}); err != nil { t.Fatalf("univ: %v", err) }
	for i := 0; i < 8; i++ {
		if err := appRepo.Create(&model.ApplicationProject{StudentID: uint(i+1), UniversityID: 1, Major:"CS", Status:"planning"}); err != nil { t.Fatalf("app: %v", err) }
	}
	return NewDashboardHandler(appSvc, service.NewUniversityService(univRepo, logger), matSvc, logger)
}

func TestDashboardStatsConcurrentRaceP801(t *testing.T) {
	h := setupDashboardHandler(t)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.GET("/stats", h.Stats)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			req := httptest.NewRequest(http.MethodGet, "/stats", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
		}()
	}
	close(start)
	wg.Wait()
}

package service

import (
	"sync"
	"testing"

	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

func TestFillAppStatsNoStaleKeysP802(t *testing.T) {
	var dst AppStats
	FillAppStats([]model.ApplicationProject{{Status: "admitted"}}, &dst)
	FillAppStats([]model.ApplicationProject{{Status: "planning"}}, &dst)
	if _, ok := dst.ByStatus["admitted"]; ok {
		t.Fatalf("stale status leaked into second fill: %+v", dst.ByStatus)
	}
}

func TestFillAppStatsResetsCountersP803(t *testing.T) {
	var dst AppStats
	FillAppStats([]model.ApplicationProject{{Status: "admitted"}}, &dst)
	FillAppStats([]model.ApplicationProject{{Status: "planning"}}, &dst)
	if dst.Admitted != 0 || dst.Applied != 0 {
		t.Fatalf("counters not reset: %+v", dst)
	}
}

func TestComputeAppStatsConcurrentRaceP807(t *testing.T) {
	projects := make([]model.ApplicationProject, 20)
	for i := range projects {
		projects[i] = model.ApplicationProject{Status: "planning"}
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 50; j++ {
				_ = ComputeAppStats(projects)
			}
		}()
	}
	close(start)
	wg.Wait()
}

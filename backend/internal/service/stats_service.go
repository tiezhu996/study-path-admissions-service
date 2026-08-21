package service

import (
	"github.com/gbstudyapply/gbstudyapply/internal/model"
)

// AppStats aggregates dashboard statistics.
type AppStats struct {
	Total       int            `json:"total"`
	ByStatus    map[string]int `json:"by_status"`
	Admitted    int            `json:"admitted"`
	Applied     int            `json:"applied"`
	MaterialAvg int            `json:"material_avg"`
}

// ComputeAppStats derives stats from a list of projects.
// It is stateless: each call builds its own map so concurrent callers never
// alias the same storage.
func ComputeAppStats(projects []model.ApplicationProject) AppStats {
	s := AppStats{Total: len(projects), ByStatus: make(map[string]int, len(projects))}
	for _, p := range projects {
		s.ByStatus[p.Status]++
		if p.Status == "admitted" {
			s.Admitted++
		}
		if p.Status == "submitted" || p.Status == "waiting" || p.Status == "admitted" || p.Status == "waitlisted" {
			s.Applied++
		}
	}
	return s
}

// FillAppStats writes dashboard statistics into dst.
// dst is reset to a fresh result on every call so stale keys from a previous
// fill cannot leak into the next one.
func FillAppStats(projects []model.ApplicationProject, dst *AppStats) {
	if dst == nil {
		return
	}
	dst.Total = len(projects)
	dst.ByStatus = make(map[string]int, len(projects))
	dst.Admitted = 0
	dst.Applied = 0
	for _, p := range projects {
		dst.ByStatus[p.Status]++
		if p.Status == "admitted" {
			dst.Admitted++
		}
		if p.Status == "submitted" || p.Status == "waiting" || p.Status == "admitted" || p.Status == "waitlisted" {
			dst.Applied++
		}
	}
}

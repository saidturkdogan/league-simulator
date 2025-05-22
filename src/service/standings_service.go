package service

import (
	"context"

	"github.com/user/league-simulator/src/model"
	"github.com/user/league-simulator/src/repository"
)

// StandingsService handles business logic for standings
type StandingsService struct {
	repo repository.StandingsRepository
}

// NewStandingsService creates a new StandingsService
func NewStandingsService(repo repository.StandingsRepository) *StandingsService {
	return &StandingsService{
		repo: repo,
	}
}

// GetCurrent retrieves the current standings
func (s *StandingsService) GetCurrent(ctx context.Context) (*model.Standings, error) {
	return s.repo.GetCurrent(ctx)
}

// Update updates the standings
func (s *StandingsService) Update(ctx context.Context, standings *model.Standings) error {
	return s.repo.Update(ctx, standings)
}

package service

import (
	"context"
	"errors"

	"github.com/user/league-simulator/src/model"
	"github.com/user/league-simulator/src/repository"
)

// MatchService handles business logic for matches
type MatchService struct {
	repo repository.MatchRepository
}

// NewMatchService creates a new MatchService
func NewMatchService(repo repository.MatchRepository) *MatchService {
	return &MatchService{
		repo: repo,
	}
}

// Create creates a new match
func (s *MatchService) Create(ctx context.Context, match *model.Match) error {
	if err := match.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, match)
}

// GetByID retrieves a match by its ID
func (s *MatchService) GetByID(ctx context.Context, id int) (*model.Match, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByWeek retrieves all matches for a specific week
func (s *MatchService) GetByWeek(ctx context.Context, week int) ([]*model.Match, error) {
	if week < 1 {
		return nil, errors.New("week must be a positive number")
	}
	return s.repo.GetByWeek(ctx, week)
}

// GetAll retrieves all matches
func (s *MatchService) GetAll(ctx context.Context) ([]*model.Match, error) {
	return s.repo.GetAll(ctx)
}

// Update updates a match
func (s *MatchService) Update(ctx context.Context, match *model.Match) error {
	if err := match.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, match)
}

// Delete removes a match
func (s *MatchService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

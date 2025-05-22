package service

import (
	"context"

	"github.com/user/league-simulator/src/model"
	"github.com/user/league-simulator/src/repository"
)

// TeamService handles business logic for teams
type TeamService struct {
	repo repository.TeamRepository
}

// NewTeamService creates a new TeamService
func NewTeamService(repo repository.TeamRepository) *TeamService {
	return &TeamService{
		repo: repo,
	}
}

// Create creates a new team
func (s *TeamService) Create(ctx context.Context, team *model.Team) error {
	if err := team.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, team)
}

// GetByID retrieves a team by its ID
func (s *TeamService) GetByID(ctx context.Context, id int) (*model.Team, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAll retrieves all teams
func (s *TeamService) GetAll(ctx context.Context) ([]*model.Team, error) {
	return s.repo.GetAll(ctx)
}

// Update updates a team
func (s *TeamService) Update(ctx context.Context, team *model.Team) error {
	if err := team.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, team)
}

// Delete removes a team
func (s *TeamService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// CreateInitialTeams creates the initial 4 teams for the league
func (s *TeamService) CreateInitialTeams(ctx context.Context) ([]*model.Team, error) {
	// Check if teams already exist
	existingTeams, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(existingTeams) > 0 {
		return existingTeams, nil
	}

	// Create 4 teams with different strengths
	teams := []*model.Team{
		{Name: "Manchester United", Strength: 85},
		{Name: "Liverpool", Strength: 88},
		{Name: "Chelsea", Strength: 82},
		{Name: "Arsenal", Strength: 80},
	}

	for _, team := range teams {
		if err := s.repo.Create(ctx, team); err != nil {
			return nil, err
		}
	}

	return teams, nil
}

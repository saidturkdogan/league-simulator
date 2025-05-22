package repository

import (
	"context"

	"github.com/user/league-simulator/src/model"
)

// TeamRepository defines the interface for team data operations
type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	GetByID(ctx context.Context, id int) (*model.Team, error)
	GetAll(ctx context.Context) ([]*model.Team, error)
	Update(ctx context.Context, team *model.Team) error
	Delete(ctx context.Context, id int) error
}

// MatchRepository defines the interface for match data operations
type MatchRepository interface {
	Create(ctx context.Context, match *model.Match) error
	GetByID(ctx context.Context, id int) (*model.Match, error)
	GetByWeek(ctx context.Context, week int) ([]*model.Match, error)
	GetAll(ctx context.Context) ([]*model.Match, error)
	Update(ctx context.Context, match *model.Match) error
	Delete(ctx context.Context, id int) error
}

// StandingsRepository defines the interface for standings data operations
type StandingsRepository interface {
	GetCurrent(ctx context.Context) (*model.Standings, error)
	Update(ctx context.Context, standings *model.Standings) error
}

// LeagueRepository defines the interface for league data operations
type LeagueRepository interface {
	Create(ctx context.Context, league *model.League) error
	GetByID(ctx context.Context, id int) (*model.League, error)
	Update(ctx context.Context, league *model.League) error
}

// Repository combines all repositories
type Repository struct {
	Team      TeamRepository
	Match     MatchRepository
	Standings StandingsRepository
	League    LeagueRepository
}

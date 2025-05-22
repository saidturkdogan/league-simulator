package repository

import (
	"database/sql"
)

// PostgresRepository implements all repository interfaces using PostgreSQL
type PostgresRepository struct {
	Team      TeamRepository
	Match     MatchRepository
	Standings StandingsRepository
	League    LeagueRepository
}

// NewPostgresRepository creates a new PostgresRepository with all implementations
func NewPostgresRepository(db *sql.DB) *Repository {
	return &Repository{
		Team:      NewPostgresTeamRepository(db),
		Match:     NewPostgresMatchRepository(db),
		Standings: NewPostgresStandingsRepository(db),
		League:    NewPostgresLeagueRepository(db),
	}
}

package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/user/league-simulator/src/model"
)

// PostgresMatchRepository implements the MatchRepository interface
type PostgresMatchRepository struct {
	db *sql.DB
}

// NewPostgresMatchRepository creates a new PostgresMatchRepository
func NewPostgresMatchRepository(db *sql.DB) *PostgresMatchRepository {
	return &PostgresMatchRepository{
		db: db,
	}
}

// Create inserts a new match into the database
func (r *PostgresMatchRepository) Create(ctx context.Context, match *model.Match) error {
	query := `
		INSERT INTO matches (home_team_id, away_team_id, home_score, away_score, week, played, played_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		match.HomeTeamID,
		match.AwayTeamID,
		match.HomeScore,
		match.AwayScore,
		match.Week,
		match.Played,
		match.PlayedAt,
	).Scan(&match.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a match by its ID
func (r *PostgresMatchRepository) GetByID(ctx context.Context, id int) (*model.Match, error) {
	query := `
		SELECT m.id, m.home_team_id, m.away_team_id, m.home_score, m.away_score, m.week, m.played, m.played_at,
			   ht.id, ht.name, ht.strength,
			   at.id, at.name, at.strength
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		WHERE m.id = $1
	`

	var match model.Match
	var homeTeam model.Team
	var awayTeam model.Team
	var playedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&match.ID,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&match.HomeScore,
		&match.AwayScore,
		&match.Week,
		&match.Played,
		&playedAt,
		&homeTeam.ID,
		&homeTeam.Name,
		&homeTeam.Strength,
		&awayTeam.ID,
		&awayTeam.Name,
		&awayTeam.Strength,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("match not found")
		}
		return nil, err
	}

	if playedAt.Valid {
		match.PlayedAt = playedAt.Time
	}

	match.HomeTeam = &homeTeam
	match.AwayTeam = &awayTeam

	return &match, nil
}

// GetByWeek retrieves all matches for a specific week
func (r *PostgresMatchRepository) GetByWeek(ctx context.Context, week int) ([]*model.Match, error) {
	query := `
		SELECT m.id, m.home_team_id, m.away_team_id, m.home_score, m.away_score, m.week, m.played, m.played_at,
			   ht.id, ht.name, ht.strength,
			   at.id, at.name, at.strength
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		WHERE m.week = $1
		ORDER BY m.id
	`

	rows, err := r.db.QueryContext(ctx, query, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*model.Match
	for rows.Next() {
		var match model.Match
		var homeTeam model.Team
		var awayTeam model.Team
		var playedAt sql.NullTime

		if err := rows.Scan(
			&match.ID,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.HomeScore,
			&match.AwayScore,
			&match.Week,
			&match.Played,
			&playedAt,
			&homeTeam.ID,
			&homeTeam.Name,
			&homeTeam.Strength,
			&awayTeam.ID,
			&awayTeam.Name,
			&awayTeam.Strength,
		); err != nil {
			return nil, err
		}

		if playedAt.Valid {
			match.PlayedAt = playedAt.Time
		}

		match.HomeTeam = &homeTeam
		match.AwayTeam = &awayTeam

		matches = append(matches, &match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// GetAll retrieves all matches
func (r *PostgresMatchRepository) GetAll(ctx context.Context) ([]*model.Match, error) {
	query := `
		SELECT m.id, m.home_team_id, m.away_team_id, m.home_score, m.away_score, m.week, m.played, m.played_at
		FROM matches m
		ORDER BY m.week, m.id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*model.Match
	for rows.Next() {
		var match model.Match
		var playedAt sql.NullTime

		if err := rows.Scan(
			&match.ID,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.HomeScore,
			&match.AwayScore,
			&match.Week,
			&match.Played,
			&playedAt,
		); err != nil {
			return nil, err
		}

		if playedAt.Valid {
			match.PlayedAt = playedAt.Time
		}

		matches = append(matches, &match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// Update updates a match
func (r *PostgresMatchRepository) Update(ctx context.Context, match *model.Match) error {
	query := `
		UPDATE matches
		SET home_team_id = $1, away_team_id = $2, home_score = $3, away_score = $4, 
			week = $5, played = $6, played_at = $7
		WHERE id = $8
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		match.HomeTeamID,
		match.AwayTeamID,
		match.HomeScore,
		match.AwayScore,
		match.Week,
		match.Played,
		match.PlayedAt,
		match.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("match not found")
	}

	return nil
}

// Delete removes a match
func (r *PostgresMatchRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM matches
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("match not found")
	}

	return nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/user/league-simulator/src/model"
)

// PostgresLeagueRepository implements the LeagueRepository interface
type PostgresLeagueRepository struct {
	db *sql.DB
}

// NewPostgresLeagueRepository creates a new PostgresLeagueRepository
func NewPostgresLeagueRepository(db *sql.DB) *PostgresLeagueRepository {
	return &PostgresLeagueRepository{
		db: db,
	}
}

// Create inserts a new league into the database
func (r *PostgresLeagueRepository) Create(ctx context.Context, league *model.League) error {
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert league
	leagueQuery := `
		INSERT INTO leagues (name, current_week, total_weeks)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err = tx.QueryRowContext(
		ctx,
		leagueQuery,
		league.Name,
		league.CurrentWeek,
		league.TotalWeeks,
	).Scan(&league.ID)
	if err != nil {
		return err
	}

	// Insert matches
	matchQuery := `
		INSERT INTO matches (home_team_id, away_team_id, week, played)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	for i := range league.Matches {
		match := league.Matches[i]
		err = tx.QueryRowContext(
			ctx,
			matchQuery,
			match.HomeTeamID,
			match.AwayTeamID,
			match.Week,
			match.Played,
		).Scan(&match.ID)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a league by its ID
func (r *PostgresLeagueRepository) GetByID(ctx context.Context, id int) (*model.League, error) {
	// Get league info
	leagueQuery := `
		SELECT id, name, current_week, total_weeks
		FROM leagues
		WHERE id = $1
	`
	league := &model.League{}
	err := r.db.QueryRowContext(ctx, leagueQuery, id).Scan(
		&league.ID,
		&league.Name,
		&league.CurrentWeek,
		&league.TotalWeeks,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("league not found")
		}
		return nil, err
	}

	// Get teams
	teamsQuery := `
		SELECT id, name, strength
		FROM teams
		ORDER BY id
	`
	teamRows, err := r.db.QueryContext(ctx, teamsQuery)
	if err != nil {
		return nil, err
	}
	defer teamRows.Close()

	var teams []*model.Team
	for teamRows.Next() {
		team := &model.Team{}
		if err := teamRows.Scan(&team.ID, &team.Name, &team.Strength); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	if err := teamRows.Err(); err != nil {
		return nil, err
	}
	league.Teams = teams

	// Get matches
	matchesQuery := `
		SELECT id, home_team_id, away_team_id, home_score, away_score, week, played, played_at
		FROM matches
		ORDER BY week, id
	`
	matchRows, err := r.db.QueryContext(ctx, matchesQuery)
	if err != nil {
		return nil, err
	}
	defer matchRows.Close()

	var matches []*model.Match
	for matchRows.Next() {
		match := &model.Match{}
		var playedAt sql.NullTime
		if err := matchRows.Scan(
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
		matches = append(matches, match)
	}
	if err := matchRows.Err(); err != nil {
		return nil, err
	}
	league.Matches = matches

	// Get standings
	standingsQuery := `
		SELECT s.team_id, t.name, s.points, s.played, s.wins, s.draws, s.losses, 
			   s.goals_for, s.goals_against, s.goals_for - s.goals_against as goal_difference
		FROM standings_history s
		JOIN teams t ON s.team_id = t.id
		WHERE s.week = $1
		ORDER BY s.points DESC, goal_difference DESC, s.goals_for DESC, t.name
	`
	standingsRows, err := r.db.QueryContext(ctx, standingsQuery, league.CurrentWeek)
	if err != nil {
		return nil, err
	}
	defer standingsRows.Close()

	standings := model.Standings{
		Week:  league.CurrentWeek,
		Teams: []model.TeamStanding{},
	}

	for standingsRows.Next() {
		var standing model.TeamStanding
		if err := standingsRows.Scan(
			&standing.TeamID,
			&standing.TeamName,
			&standing.Points,
			&standing.Played,
			&standing.Wins,
			&standing.Draws,
			&standing.Losses,
			&standing.GoalsFor,
			&standing.GoalsAgainst,
			&standing.GoalDifference,
		); err != nil {
			return nil, err
		}
		standings.Teams = append(standings.Teams, standing)
	}
	if err := standingsRows.Err(); err != nil {
		return nil, err
	}
	league.Standings = standings

	return league, nil
}

// Update updates a league
func (r *PostgresLeagueRepository) Update(ctx context.Context, league *model.League) error {
	query := `
		UPDATE leagues
		SET name = $1, current_week = $2, total_weeks = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		league.Name,
		league.CurrentWeek,
		league.TotalWeeks,
		league.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("league not found")
	}

	return nil
}

package repository

import (
	"context"
	"database/sql"

	"github.com/user/league-simulator/src/model"
)

// PostgresStandingsRepository implements the StandingsRepository interface
type PostgresStandingsRepository struct {
	db *sql.DB
}

// NewPostgresStandingsRepository creates a new PostgresStandingsRepository
func NewPostgresStandingsRepository(db *sql.DB) *PostgresStandingsRepository {
	return &PostgresStandingsRepository{
		db: db,
	}
}

// GetCurrent retrieves the current standings
func (r *PostgresStandingsRepository) GetCurrent(ctx context.Context) (*model.Standings, error) {
	// First get the current week
	weekQuery := `
		SELECT COALESCE(MAX(week), 0) FROM standings_history
	`
	var week int
	err := r.db.QueryRowContext(ctx, weekQuery).Scan(&week)
	if err != nil {
		return nil, err
	}

	// If no standings exist yet, return empty standings
	if week == 0 {
		// Get all teams to create empty standings
		teamsQuery := `
			SELECT id, name FROM teams ORDER BY id
		`
		rows, err := r.db.QueryContext(ctx, teamsQuery)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		standings := &model.Standings{
			Week:  0,
			Teams: []model.TeamStanding{},
		}

		for rows.Next() {
			var teamID int
			var teamName string
			if err := rows.Scan(&teamID, &teamName); err != nil {
				return nil, err
			}

			standings.Teams = append(standings.Teams, model.TeamStanding{
				TeamID:   teamID,
				TeamName: teamName,
			})
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		return standings, nil
	}

	// Get standings for the current week
	query := `
		SELECT s.team_id, t.name, s.points, s.played, s.wins, s.draws, s.losses, 
			   s.goals_for, s.goals_against, s.goals_for - s.goals_against as goal_difference
		FROM standings_history s
		JOIN teams t ON s.team_id = t.id
		WHERE s.week = $1
		ORDER BY s.points DESC, goal_difference DESC, s.goals_for DESC, t.name
	`

	rows, err := r.db.QueryContext(ctx, query, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	standings := &model.Standings{
		Week:  week,
		Teams: []model.TeamStanding{},
	}

	for rows.Next() {
		var standing model.TeamStanding
		if err := rows.Scan(
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return standings, nil
}

// Update updates the standings
func (r *PostgresStandingsRepository) Update(ctx context.Context, standings *model.Standings) error {
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert or update standings for each team
	for _, team := range standings.Teams {
		query := `
			INSERT INTO standings_history (
				team_id, week, points, played, wins, draws, losses, goals_for, goals_against
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9
			)
			ON CONFLICT (team_id, week) DO UPDATE SET
				points = $3, played = $4, wins = $5, draws = $6, losses = $7, 
				goals_for = $8, goals_against = $9
		`

		_, err := tx.ExecContext(
			ctx,
			query,
			team.TeamID,
			standings.Week,
			team.Points,
			team.Played,
			team.Wins,
			team.Draws,
			team.Losses,
			team.GoalsFor,
			team.GoalsAgainst,
		)
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

package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/user/league-simulator/src/model"
)

// PostgresTeamRepository implements the TeamRepository interface
type PostgresTeamRepository struct {
	db *sql.DB
}

// NewPostgresTeamRepository creates a new PostgresTeamRepository
func NewPostgresTeamRepository(db *sql.DB) *PostgresTeamRepository {
	return &PostgresTeamRepository{
		db: db,
	}
}

// Create inserts a new team into the database
func (r *PostgresTeamRepository) Create(ctx context.Context, team *model.Team) error {
	query := `
		INSERT INTO teams (name, strength)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, team.Name, team.Strength).Scan(&team.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a team by its ID
func (r *PostgresTeamRepository) GetByID(ctx context.Context, id int) (*model.Team, error) {
	query := `
		SELECT id, name, strength
		FROM teams
		WHERE id = $1
	`

	team := &model.Team{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&team.ID, &team.Name, &team.Strength)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("team not found")
		}
		return nil, err
	}

	return team, nil
}

// GetAll retrieves all teams
func (r *PostgresTeamRepository) GetAll(ctx context.Context) ([]*model.Team, error) {
	query := `
		SELECT id, name, strength
		FROM teams
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*model.Team
	for rows.Next() {
		team := &model.Team{}
		if err := rows.Scan(&team.ID, &team.Name, &team.Strength); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

// Update updates a team
func (r *PostgresTeamRepository) Update(ctx context.Context, team *model.Team) error {
	query := `
		UPDATE teams
		SET name = $1, strength = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, team.Name, team.Strength, team.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("team not found")
	}

	return nil
}

// Delete removes a team
func (r *PostgresTeamRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM teams
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
		return errors.New("team not found")
	}

	return nil
}

package model

import (
	"errors"
)

// Team represents a football team in the league
type Team struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Strength int    `json:"strength"` // 1-100 scale representing team's strength
}

// Validate checks if the team data is valid
func (t *Team) Validate() error {
	if t.Name == "" {
		return errors.New("team name cannot be empty")
	}

	if t.Strength < 1 || t.Strength > 100 {
		return errors.New("team strength must be between 1 and 100")
	}

	return nil
}

package model

import (
	"errors"
	"time"
)

// Match represents a football match between two teams
type Match struct {
	ID         int       `json:"id"`
	HomeTeamID int       `json:"home_team_id"`
	AwayTeamID int       `json:"away_team_id"`
	HomeTeam   *Team     `json:"home_team,omitempty"`
	AwayTeam   *Team     `json:"away_team,omitempty"`
	HomeScore  int       `json:"home_score"`
	AwayScore  int       `json:"away_score"`
	Week       int       `json:"week"`
	Played     bool      `json:"played"`
	PlayedAt   time.Time `json:"played_at,omitempty"`
}

// Validate checks if the match data is valid
func (m *Match) Validate() error {
	if m.HomeTeamID == m.AwayTeamID {
		return errors.New("home team and away team cannot be the same")
	}

	if m.Week < 1 {
		return errors.New("week must be a positive number")
	}

	if m.Played {
		if m.HomeScore < 0 || m.AwayScore < 0 {
			return errors.New("scores cannot be negative")
		}
	}

	return nil
}

// Result returns the match result from the perspective of the home team
func (m *Match) Result() string {
	if !m.Played {
		return "Not Played"
	}

	if m.HomeScore > m.AwayScore {
		return "Win"
	} else if m.HomeScore < m.AwayScore {
		return "Loss"
	}
	return "Draw"
}

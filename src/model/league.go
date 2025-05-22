package model

import (
	"errors"
	"math/rand"
	"time"
)

// League represents a football league
type League struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Teams     []*Team  `json:"teams"`
	Matches   []*Match `json:"matches,omitempty"`
	Standings Standings `json:"standings"`
	CurrentWeek int     `json:"current_week"`
	TotalWeeks  int     `json:"total_weeks"`
}

// NewLeague creates a new league with the given teams
func NewLeague(name string, teams []*Team) (*League, error) {
	if len(teams) < 2 {
		return nil, errors.New("league must have at least 2 teams")
	}

	// Calculate total weeks based on round-robin tournament
	totalWeeks := len(teams) - 1

	league := &League{
		Name:       name,
		Teams:      teams,
		CurrentWeek: 0,
		TotalWeeks:  totalWeeks,
		Standings: Standings{
			Teams: make([]TeamStanding, len(teams)),
			Week:  0,
		},
	}

	// Initialize standings
	for i, team := range teams {
		league.Standings.Teams[i] = TeamStanding{
			TeamID:   team.ID,
			TeamName: team.Name,
		}
	}

	// Generate match schedule
	league.generateSchedule()

	return league, nil
}

// generateSchedule creates a round-robin tournament schedule
func (l *League) generateSchedule() {
	numTeams := len(l.Teams)
	
	// For odd number of teams, add a dummy team
	if numTeams%2 != 0 {
		numTeams++
	}
	
	// Generate matches for each week
	for week := 1; week <= l.TotalWeeks; week++ {
		for i := 0; i < numTeams/2; i++ {
			// Skip matches involving the dummy team
			if i == 0 && numTeams > len(l.Teams) {
				continue
			}
			
			// Calculate indices for home and away teams
			home := (week + i) % (numTeams - 1)
			away := (numTeams - 1 - i + week) % (numTeams - 1)
			
			// Last team stays fixed
			if i == 0 {
				away = numTeams - 1
			}
			
			// Skip if either team is the dummy team
			if home >= len(l.Teams) || away >= len(l.Teams) {
				continue
			}
			
			// Create the match
			match := &Match{
				HomeTeamID: l.Teams[home].ID,
				AwayTeamID: l.Teams[away].ID,
				Week:       week,
				Played:     false,
			}
			
			l.Matches = append(l.Matches, match)
		}
	}
}

// SimulateWeek simulates all matches for the current week
func (l *League) SimulateWeek() error {
	if l.CurrentWeek >= l.TotalWeeks {
		return errors.New("all weeks have been played")
	}
	
	l.CurrentWeek++
	
	// Find matches for the current week
	for _, match := range l.Matches {
		if match.Week == l.CurrentWeek && !match.Played {
			l.SimulateMatch(match)
			l.Standings.UpdateStandings(match)
		}
	}
	
	l.Standings.Week = l.CurrentWeek
	
	return nil
}

// SimulateMatch simulates a single match based on team strengths
func (l *League) SimulateMatch(match *Match) {
	// Find the teams
	var homeTeam, awayTeam *Team
	for _, team := range l.Teams {
		if team.ID == match.HomeTeamID {
			homeTeam = team
		}
		if team.ID == match.AwayTeamID {
			awayTeam = team
		}
	}
	
	if homeTeam == nil || awayTeam == nil {
		return
	}
	
	// Home advantage factor (1.2x)
	homeAdvantage := 1.2
	
	// Calculate effective strengths
	homeStrength := float64(homeTeam.Strength) * homeAdvantage
	awayStrength := float64(awayTeam.Strength)
	
	// Random factor (0.7 to 1.3)
	// Use a local random generator (Go 1.20+ recommendation)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	homeRandom := 0.7 + r.Float64()*0.6
	awayRandom := 0.7 + r.Float64()*0.6
	
	// Calculate scores based on strengths and randomness
	homeScoreFactor := homeStrength * homeRandom / 25.0
	awayScoreFactor := awayStrength * awayRandom / 30.0
	
	// Convert to integer scores (0-5 range is common in football)
	match.HomeScore = min(int(homeScoreFactor), 5)
	match.AwayScore = min(int(awayScoreFactor), 5)
	
	match.Played = true
	match.PlayedAt = time.Now()
}

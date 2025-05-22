package model

// TeamStanding represents a team's position in the league standings
type TeamStanding struct {
	TeamID        int    `json:"team_id"`
	TeamName      string `json:"team_name"`
	Points        int    `json:"points"`
	Played        int    `json:"played"`
	Wins          int    `json:"wins"`
	Draws         int    `json:"draws"`
	Losses        int    `json:"losses"`
	GoalsFor      int    `json:"goals_for"`
	GoalsAgainst  int    `json:"goals_against"`
	GoalDifference int    `json:"goal_difference"`
}

// Standings represents the league standings
type Standings struct {
	Teams []TeamStanding `json:"teams"`
	Week  int            `json:"week"`
}

// UpdateStandings updates the standings based on a match result
func (s *Standings) UpdateStandings(match *Match) {
	if !match.Played {
		return
	}

	// Find home team in standings
	var homeTeamIndex, awayTeamIndex = -1, -1
	for i, team := range s.Teams {
		if team.TeamID == match.HomeTeamID {
			homeTeamIndex = i
		}
		if team.TeamID == match.AwayTeamID {
			awayTeamIndex = i
		}
	}

	// Update home team stats
	if homeTeamIndex >= 0 {
		s.Teams[homeTeamIndex].Played++
		s.Teams[homeTeamIndex].GoalsFor += match.HomeScore
		s.Teams[homeTeamIndex].GoalsAgainst += match.AwayScore
		s.Teams[homeTeamIndex].GoalDifference = s.Teams[homeTeamIndex].GoalsFor - s.Teams[homeTeamIndex].GoalsAgainst

		if match.HomeScore > match.AwayScore {
			s.Teams[homeTeamIndex].Wins++
			s.Teams[homeTeamIndex].Points += 3
		} else if match.HomeScore == match.AwayScore {
			s.Teams[homeTeamIndex].Draws++
			s.Teams[homeTeamIndex].Points += 1
		} else {
			s.Teams[homeTeamIndex].Losses++
		}
	}

	// Update away team stats
	if awayTeamIndex >= 0 {
		s.Teams[awayTeamIndex].Played++
		s.Teams[awayTeamIndex].GoalsFor += match.AwayScore
		s.Teams[awayTeamIndex].GoalsAgainst += match.HomeScore
		s.Teams[awayTeamIndex].GoalDifference = s.Teams[awayTeamIndex].GoalsFor - s.Teams[awayTeamIndex].GoalsAgainst

		if match.AwayScore > match.HomeScore {
			s.Teams[awayTeamIndex].Wins++
			s.Teams[awayTeamIndex].Points += 3
		} else if match.AwayScore == match.HomeScore {
			s.Teams[awayTeamIndex].Draws++
			s.Teams[awayTeamIndex].Points += 1
		} else {
			s.Teams[awayTeamIndex].Losses++
		}
	}
}

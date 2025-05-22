package service

import (
	"context"
	"sort"

	"github.com/user/league-simulator/src/model"
	"github.com/user/league-simulator/src/repository"
)

// PredictionService handles prediction logic
type PredictionService struct {
	leagueRepo repository.LeagueRepository
	teamRepo   repository.TeamRepository
	matchRepo  repository.MatchRepository
}

// NewPredictionService creates a new PredictionService
func NewPredictionService(
	leagueRepo repository.LeagueRepository,
	teamRepo repository.TeamRepository,
	matchRepo repository.MatchRepository,
) *PredictionService {
	return &PredictionService{
		leagueRepo: leagueRepo,
		teamRepo:   teamRepo,
		matchRepo:  matchRepo,
	}
}

// PredictFinalStandings predicts the final standings after all weeks
func (s *PredictionService) PredictFinalStandings(ctx context.Context, leagueID int) (*model.Standings, error) {
	// Get the league
	league, err := s.leagueRepo.GetByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	// If all weeks have been played, return the current standings
	if league.CurrentWeek >= league.TotalWeeks {
		return &league.Standings, nil
	}

	// Create a copy of the current standings
	predictedStandings := model.Standings{
		Week:  league.TotalWeeks,
		Teams: make([]model.TeamStanding, len(league.Standings.Teams)),
	}
	copy(predictedStandings.Teams, league.Standings.Teams)

	// Find remaining matches
	var remainingMatches []*model.Match
	for _, match := range league.Matches {
		if !match.Played {
			remainingMatches = append(remainingMatches, match)
		}
	}

	// Simulate remaining matches
	// No need to seed the random generator in Go 1.20+
	for _, match := range remainingMatches {
		// Find the teams
		var homeTeam, awayTeam *model.Team
		for _, team := range league.Teams {
			if team.ID == match.HomeTeamID {
				homeTeam = team
			}
			if team.ID == match.AwayTeamID {
				awayTeam = team
			}
		}

		if homeTeam == nil || awayTeam == nil {
			continue
		}

		// Create a copy of the match for simulation
		simulatedMatch := &model.Match{
			HomeTeamID: match.HomeTeamID,
			AwayTeamID: match.AwayTeamID,
			HomeTeam:   homeTeam,
			AwayTeam:   awayTeam,
			Week:       match.Week,
		}

		// Simulate the match
		// Create a temporary league for simulation
		tempLeague := &model.League{
			Teams: league.Teams,
		}
		tempLeague.SimulateMatch(simulatedMatch)

		// Update the predicted standings
		predictedStandings.UpdateStandings(simulatedMatch)
	}

	// Sort the standings by points, goal difference, etc.
	sort.Slice(predictedStandings.Teams, func(i, j int) bool {
		if predictedStandings.Teams[i].Points != predictedStandings.Teams[j].Points {
			return predictedStandings.Teams[i].Points > predictedStandings.Teams[j].Points
		}
		if predictedStandings.Teams[i].GoalDifference != predictedStandings.Teams[j].GoalDifference {
			return predictedStandings.Teams[i].GoalDifference > predictedStandings.Teams[j].GoalDifference
		}
		if predictedStandings.Teams[i].GoalsFor != predictedStandings.Teams[j].GoalsFor {
			return predictedStandings.Teams[i].GoalsFor > predictedStandings.Teams[j].GoalsFor
		}
		return predictedStandings.Teams[i].TeamName < predictedStandings.Teams[j].TeamName
	})

	return &predictedStandings, nil
}

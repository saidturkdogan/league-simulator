package service

import (
	"context"
	"errors"
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
	s.sortStandings(&predictedStandings)

	return &predictedStandings, nil
}

// GetPredictionWithConfidence returns predictions with confidence levels after week 4
func (s *PredictionService) GetPredictionWithConfidence(ctx context.Context, leagueID int) (*model.PredictionResult, error) {
	// Get the league
	league, err := s.leagueRepo.GetByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	// Predictions are only available after week 4 as per requirements
	if league.CurrentWeek < 4 {
		return nil, errors.New("tahminler sadece 4. hafta sonrası kullanılabilir")
	}

	// If all weeks have been played, return final results
	if league.CurrentWeek >= league.TotalWeeks {
		finalStandings := league.Standings
		// Sort standings
		s.sortStandings(&finalStandings)
		
		return &model.PredictionResult{
			CurrentWeek:    league.CurrentWeek,
			TotalWeeks:     league.TotalWeeks,
			PredictionType: "Final Results",
			Standings:      &finalStandings,
			Confidence:     100.0, // 100% confidence for final results
		}, nil
	}

	// Run multiple simulations for better prediction accuracy
	simulations := 100
	teamPredictions := make(map[int]*model.TeamPrediction)

	// Initialize predictions
	for _, team := range league.Teams {
		teamPredictions[team.ID] = &model.TeamPrediction{
			TeamID:        team.ID,
			TeamName:      team.Name,
			CurrentPoints: 0,
			PositionCounts: make([]int, len(league.Teams)),
		}
	}

	// Get current points from standings
	for _, standing := range league.Standings.Teams {
		if pred, exists := teamPredictions[standing.TeamID]; exists {
			pred.CurrentPoints = standing.Points
		}
	}

	// Run multiple simulations
	for sim := 0; sim < simulations; sim++ {
		predictedStandings, err := s.PredictFinalStandings(ctx, leagueID)
		if err != nil {
			continue
		}

		// Record positions from this simulation
		for pos, standing := range predictedStandings.Teams {
			if pred, exists := teamPredictions[standing.TeamID]; exists {
				pred.PositionCounts[pos]++
				pred.PredictedPoints += standing.Points
			}
		}
	}

	// Calculate averages and probabilities
	for _, pred := range teamPredictions {
		pred.PredictedPoints /= simulations
		pred.ChampionshipProbability = float64(pred.PositionCounts[0]) / float64(simulations) * 100
		pred.TopThreeProbability = float64(pred.PositionCounts[0]+pred.PositionCounts[1]+pred.PositionCounts[2]) / float64(simulations) * 100
		pred.RelegationProbability = float64(pred.PositionCounts[len(pred.PositionCounts)-1]) / float64(simulations) * 100
		
		// Calculate most likely position
		maxCount := 0
		for pos, count := range pred.PositionCounts {
			if count > maxCount {
				maxCount = count
				pred.MostLikelyPosition = pos + 1 // 1-indexed
			}
		}
	}

	// Calculate overall confidence based on matches played
	matchesPlayed := league.CurrentWeek
	totalMatches := league.TotalWeeks
	confidence := float64(matchesPlayed) / float64(totalMatches) * 100

	result := &model.PredictionResult{
		CurrentWeek:       league.CurrentWeek,
		TotalWeeks:        league.TotalWeeks,
		PredictionType:    "Statistical Prediction",
		TeamPredictions:   make([]*model.TeamPrediction, 0),
		Confidence:        confidence,
	}

	// Convert map to slice and sort by predicted points
	for _, pred := range teamPredictions {
		result.TeamPredictions = append(result.TeamPredictions, pred)
	}

	sort.Slice(result.TeamPredictions, func(i, j int) bool {
		return result.TeamPredictions[i].PredictedPoints > result.TeamPredictions[j].PredictedPoints
	})

	// Get one final prediction for the standings field
	finalPrediction, err := s.PredictFinalStandings(ctx, leagueID)
	if err != nil {
		return nil, err
	}
	result.Standings = finalPrediction

	return result, nil
}

// Helper function to sort standings
func (s *PredictionService) sortStandings(standings *model.Standings) {
	sort.Slice(standings.Teams, func(i, j int) bool {
		if standings.Teams[i].Points != standings.Teams[j].Points {
			return standings.Teams[i].Points > standings.Teams[j].Points
		}
		if standings.Teams[i].GoalDifference != standings.Teams[j].GoalDifference {
			return standings.Teams[i].GoalDifference > standings.Teams[j].GoalDifference
		}
		if standings.Teams[i].GoalsFor != standings.Teams[j].GoalsFor {
			return standings.Teams[i].GoalsFor > standings.Teams[j].GoalsFor
		}
		return standings.Teams[i].TeamName < standings.Teams[j].TeamName
	})
}

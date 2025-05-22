package service

import (
	"github.com/user/league-simulator/src/repository"
)

// Service combines all services
type Service struct {
	Team       *TeamService
	Match      *MatchService
	Standings  *StandingsService
	League     *LeagueService
	Prediction *PredictionService
}

// NewService creates a new Service with all service implementations
func NewService(repo *repository.Repository) *Service {
	return &Service{
		Team:       NewTeamService(repo.Team),
		Match:      NewMatchService(repo.Match),
		Standings:  NewStandingsService(repo.Standings),
		League:     NewLeagueService(repo.League, repo.Team, repo.Match, repo.Standings),
		Prediction: NewPredictionService(repo.League, repo.Team, repo.Match),
	}
}

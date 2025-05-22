package service

import (
	"context"
	"errors"

	"github.com/user/league-simulator/src/model"
	"github.com/user/league-simulator/src/repository"
)

// LeagueService handles business logic for leagues
type LeagueService struct {
	leagueRepo    repository.LeagueRepository
	teamRepo      repository.TeamRepository
	matchRepo     repository.MatchRepository
	standingsRepo repository.StandingsRepository
}

// NewLeagueService creates a new LeagueService
func NewLeagueService(
	leagueRepo repository.LeagueRepository,
	teamRepo repository.TeamRepository,
	matchRepo repository.MatchRepository,
	standingsRepo repository.StandingsRepository,
) *LeagueService {
	return &LeagueService{
		leagueRepo:    leagueRepo,
		teamRepo:      teamRepo,
		matchRepo:     matchRepo,
		standingsRepo: standingsRepo,
	}
}

// Create creates a new league
func (s *LeagueService) Create(ctx context.Context, name string) (*model.League, error) {
	// Get all teams
	teams, err := s.teamRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(teams) < 2 {
		return nil, errors.New("at least 2 teams are required to create a league")
	}

	// Create a new league
	league, err := model.NewLeague(name, teams)
	if err != nil {
		return nil, err
	}

	// Save the league
	if err := s.leagueRepo.Create(ctx, league); err != nil {
		return nil, err
	}

	return league, nil
}

// GetByID retrieves a league by its ID
func (s *LeagueService) GetByID(ctx context.Context, id int) (*model.League, error) {
	return s.leagueRepo.GetByID(ctx, id)
}

// SimulateWeek simulates all matches for the current week
func (s *LeagueService) SimulateWeek(ctx context.Context, leagueID int) (*model.Standings, error) {
	// Get the league
	league, err := s.leagueRepo.GetByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	if league.CurrentWeek >= league.TotalWeeks {
		return nil, errors.New("all weeks have been played")
	}

	// Increment the current week
	league.CurrentWeek++

	// Find matches for the current week
	var weekMatches []*model.Match
	for _, match := range league.Matches {
		if match.Week == league.CurrentWeek && !match.Played {
			weekMatches = append(weekMatches, match)
		}
	}

	// Simulate each match
	for _, match := range weekMatches {
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

		// Simulate the match
		match.HomeTeam = homeTeam
		match.AwayTeam = awayTeam
		league.SimulateMatch(match)

		// Update the match in the database
		if err := s.matchRepo.Update(ctx, match); err != nil {
			return nil, err
		}

		// Update standings
		league.Standings.UpdateStandings(match)
	}

	// Update the league
	if err := s.leagueRepo.Update(ctx, league); err != nil {
		return nil, err
	}

	// Update the standings
	league.Standings.Week = league.CurrentWeek
	if err := s.standingsRepo.Update(ctx, &league.Standings); err != nil {
		return nil, err
	}

	return &league.Standings, nil
}

// GetCurrentStandings retrieves the current standings
func (s *LeagueService) GetCurrentStandings(ctx context.Context, leagueID int) (*model.Standings, error) {
	// Get the league
	league, err := s.leagueRepo.GetByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	return &league.Standings, nil
}

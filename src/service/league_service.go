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

// SimulateAllRemainingWeeks - Kalan tüm haftaları otomatik simüle eder
func (s *LeagueService) SimulateAllRemainingWeeks(ctx context.Context, leagueID int) (*model.LeagueSimulationResult, error) {
	// Liga bilgilerini al
	league, err := s.leagueRepo.GetByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	if league.CurrentWeek >= league.TotalWeeks {
		return nil, errors.New("tüm haftalar zaten oynanmış")
	}

	result := &model.LeagueSimulationResult{
		LeagueID:      leagueID,
		StartingWeek:  league.CurrentWeek + 1,
		EndingWeek:    league.TotalWeeks,
		WeeklyResults: make([]*model.WeeklyResult, 0),
	}

	// Kalan haftaları tek tek simüle et
	for league.CurrentWeek < league.TotalWeeks {
		// Haftayı artır
		league.CurrentWeek++

		// Bu haftanın maçlarını bul
		var weekMatches []*model.Match
		for _, match := range league.Matches {
			if match.Week == league.CurrentWeek && !match.Played {
				weekMatches = append(weekMatches, match)
			}
		}

		weekResult := &model.WeeklyResult{
			Week:            league.CurrentWeek,
			Matches:         make([]*model.MatchResult, 0),
			StandingsBefore: s.copyStandings(&league.Standings),
		}

		// Her maçı simüle et
		for _, match := range weekMatches {
			// Takımları bul
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

			// Maçı simüle et
			match.HomeTeam = homeTeam
			match.AwayTeam = awayTeam
			league.SimulateMatch(match)

			// Maç sonucunu kaydet
			matchResult := &model.MatchResult{
				MatchID:   match.ID,
				HomeTeam:  homeTeam.Name,
				AwayTeam:  awayTeam.Name,
				HomeScore: match.HomeScore,
				AwayScore: match.AwayScore,
				Result:    match.Result(),
				PlayedAt:  match.PlayedAt,
			}
			weekResult.Matches = append(weekResult.Matches, matchResult)

			// Veritabanını güncelle
			if err := s.matchRepo.Update(ctx, match); err != nil {
				return nil, err
			}

			// Puan tablosunu güncelle
			league.Standings.UpdateStandings(match)
		}

		// Liga ve puan tablosunu kaydet
		if err := s.leagueRepo.Update(ctx, league); err != nil {
			return nil, err
		}

		league.Standings.Week = league.CurrentWeek
		if err := s.standingsRepo.Update(ctx, &league.Standings); err != nil {
			return nil, err
		}

		weekResult.StandingsAfter = s.copyStandings(&league.Standings)
		result.WeeklyResults = append(result.WeeklyResults, weekResult)
	}

	result.FinalStandings = &league.Standings
	return result, nil
}

// EditMatchResult - Maç sonucunu düzenler ve puan tablosunu yeniden hesaplar
func (s *LeagueService) EditMatchResult(ctx context.Context, matchID int, homeScore, awayScore int) (*model.Standings, error) {
	// Maçı bul
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, err
	}

	if !match.Played {
		return nil, errors.New("oynanmamış maçın sonucu düzenlenemez")
	}

	if homeScore < 0 || awayScore < 0 {
		return nil, errors.New("skorlar negatif olamaz")
	}

	// Skorları güncelle
	match.HomeScore = homeScore
	match.AwayScore = awayScore

	// Veritabanını güncelle
	if err := s.matchRepo.Update(ctx, match); err != nil {
		return nil, err
	}

	// Puan tablosunu yeniden hesapla
	return s.recalculateStandings(ctx, match)
}

// GetWeeklyMatches - Belirli bir haftanın maçlarını getir
func (s *LeagueService) GetWeeklyMatches(ctx context.Context, leagueID, week int) ([]*model.Match, error) {
	league, err := s.leagueRepo.GetByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	if week < 1 || week > league.TotalWeeks {
		return nil, errors.New("geçersiz hafta numarası")
	}

	var weekMatches []*model.Match
	for _, match := range league.Matches {
		if match.Week == week {
			// Takım bilgilerini ekle
			for _, team := range league.Teams {
				if team.ID == match.HomeTeamID {
					match.HomeTeam = team
				}
				if team.ID == match.AwayTeamID {
					match.AwayTeam = team
				}
			}
			weekMatches = append(weekMatches, match)
		}
	}

	return weekMatches, nil
}

// Helper fonksiyonlar
func (s *LeagueService) copyStandings(standings *model.Standings) *model.Standings {
	copy := &model.Standings{
		Week:  standings.Week,
		Teams: make([]model.TeamStanding, len(standings.Teams)),
	}
	for i, team := range standings.Teams {
		copy.Teams[i] = team
	}
	return copy
}

func (s *LeagueService) recalculateStandings(ctx context.Context, editedMatch *model.Match) (*model.Standings, error) {
	// Tüm takımları al
	teams, err := s.teamRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Temiz puan tablosu başlat
	standings := &model.Standings{
		Week:  editedMatch.Week,
		Teams: make([]model.TeamStanding, len(teams)),
	}

	for i, team := range teams {
		standings.Teams[i] = model.TeamStanding{
			TeamID:   team.ID,
			TeamName: team.Name,
		}
	}

	// Tüm oynanmış maçları al ve puan tablosunu yeniden hesapla
	allMatches, err := s.matchRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, match := range allMatches {
		if match.Played && match.Week <= editedMatch.Week {
			standings.UpdateStandings(match)
		}
	}

	// Veritabanını güncelle
	if err := s.standingsRepo.Update(ctx, standings); err != nil {
		return nil, err
	}

	return standings, nil
}

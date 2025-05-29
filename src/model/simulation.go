package model

import "time"

// LeagueSimulationResult - Tüm liganın simülasyon sonuçlarını tutar
type LeagueSimulationResult struct {
	LeagueID        int              `json:"league_id"`        // Liga ID'si
	StartingWeek    int              `json:"starting_week"`    // Başlangıç haftası
	EndingWeek      int              `json:"ending_week"`      // Bitiş haftası
	WeeklyResults   []*WeeklyResult  `json:"weekly_results"`   // Haftalık sonuçlar
	FinalStandings  *Standings       `json:"final_standings"`  // Final puan tablosu
}

// WeeklyResult - Bir haftanın sonuçlarını tutar
type WeeklyResult struct {
	Week            int            `json:"week"`              // Hangi hafta
	Matches         []*MatchResult `json:"matches"`           // O haftanın maçları
	StandingsBefore *Standings     `json:"standings_before"`  // Hafta öncesi puan durumu
	StandingsAfter  *Standings     `json:"standings_after"`   // Hafta sonrası puan durumu
}

// MatchResult - Bir maçın sonucunu tutar
type MatchResult struct {
	MatchID   int       `json:"match_id"`    // Maç ID'si
	HomeTeam  string    `json:"home_team"`   // Ev sahibi takım
	AwayTeam  string    `json:"away_team"`   // Deplasman takımı
	HomeScore int       `json:"home_score"`  // Ev sahibi skoru
	AwayScore int       `json:"away_score"`  // Deplasman skoru
	Result    string    `json:"result"`      // Sonuç (Win/Draw/Loss)
	PlayedAt  time.Time `json:"played_at"`   // Oynanma zamanı
}

// EditMatchRequest - Maç sonucu düzenleme talebi
type EditMatchRequest struct {
	HomeScore int `json:"home_score"` // Yeni ev sahibi skoru
	AwayScore int `json:"away_score"` // Yeni deplasman skoru
} 
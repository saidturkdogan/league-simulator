package model

// PredictionResult - 4. hafta sonrası tahmin sonuçları
type PredictionResult struct {
	CurrentWeek       int               `json:"current_week"`      // Şu anki hafta
	TotalWeeks        int               `json:"total_weeks"`       // Toplam hafta sayısı
	PredictionType    string            `json:"prediction_type"`   // Tahmin türü
	Standings         *Standings        `json:"predicted_standings"` // Tahmini puan tablosu
	TeamPredictions   []*TeamPrediction `json:"team_predictions"`  // Takım bazlı tahminler
	Confidence        float64           `json:"confidence_percentage"` // Güven yüzdesi
}

// TeamPrediction - Bir takım için detaylı tahmin
type TeamPrediction struct {
	TeamID                   int     `json:"team_id"`                   // Takım ID'si
	TeamName                 string  `json:"team_name"`                 // Takım adı
	CurrentPoints            int     `json:"current_points"`            // Şu anki puanı
	PredictedPoints          int     `json:"predicted_points"`          // Tahmini final puanı
	MostLikelyPosition       int     `json:"most_likely_position"`      // En olası sıralaması
	ChampionshipProbability  float64 `json:"championship_probability"`  // Şampiyonluk olasılığı
	TopThreeProbability      float64 `json:"top_three_probability"`     // İlk 3'e girme olasılığı
	RelegationProbability    float64 `json:"relegation_probability"`    // Küme düşme olasılığı
	PositionCounts           []int   `json:"-"`                         // Hesaplama için kullanılır
} 
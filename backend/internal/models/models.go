package models

import "time"

type Record struct {
	ID           uint      `gorm:"primaryKey"` // Auto-increment ID
	AssetType    string    `json:"asset_type"` // Crypto or ETF
	Timestamp    time.Time `json:"timestamp"`
	BidAskSpread float64   `json:"bid_ask_spread"` // Difference between ask and bid prices
	Volume       float64   `json:"volume"`         // Trading volume
	BidPrice     float64   `json:"bid_price"`      // High price (useful for trend analysis)
}

type LiquidityReport struct {
	AssetType                  string   `json:"asset_type"`
	TotalRecords               int      `json:"total_records"`
	HighRiskCount              int      `json:"high_risk_count"`
	ModerateRiskCount          int      `json:"moderate_risk_count"`
	CurrentWarnings            []string `json:"current_warnings"`
	PredictedWarnings          []string `json:"predicted_warnings"`
	CurrentHighRiskCount       int      `json:"current_high_risk_count"`
	PredictedHighRiskCount     int      `json:"predicted_high_risk_count"`
	CurrentModerateRiskCount   int      `json:"current_moderate_risk_count"`
	PredictedModerateRiskCount int      `json:"predicted_moderate_risk_count"`
}

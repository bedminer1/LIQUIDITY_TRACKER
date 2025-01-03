package models

type Record struct {
	ID           uint    `gorm:"primaryKey"` // Auto-increment ID
	AssetType    string  `json:"asset_type"` // Crypto or ETF
	Timestamp    string  `json:"timestamp"`
	BidAskSpread float64 `json:"bid_ask_spread"` // Difference between ask and bid prices
	Volume       float64 `json:"volume"`         // Trading volume
	BidPrice     float64 `json:"bid_price"`      // High price (useful for trend analysis)
}

type LiquidityReport struct {
	AssetType         string   `json:"asset_type"`
	TotalRecords      int      `json:"total_records"`
	HighRiskCount     int      `json:"high_risk_count"`
	ModerateRiskCount int      `json:"moderate_risk_count"`
	LowRiskCount      int      `json:"low_risk_count"`
	Warnings          []string `json:"warnings"`
}

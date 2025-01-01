package models

type Record struct {
	ID           uint    `gorm:"primaryKey"` // Auto-increment ID
	AssetType    string  `json:"asset_type"` // Crypto or ETF
	Timestamp    string  `json:"timestamp"`
	BidAskSpread float64 `json:"bid_ask_spread"` // Difference between ask and bid prices
	Volume       float64 `json:"volume"`         // Trading volume
	High         float64 `json:"high"`           // High price (useful for trend analysis)
	Low          float64 `json:"low"`            // Low price (useful for trend analysis)
}

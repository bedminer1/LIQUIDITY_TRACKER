package models

import "time"

type Record struct {
	ID           uint      `gorm:"primaryKey" json:"-"` // Auto-increment ID
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

type TransactionRecord struct {
	ID                           uint    `gorm:"primaryKey" json:"id,omitempty"` // Primary key
	DistanceFromHome             float64 `json:"distance_from_home"`             // Distance from home
	DistanceFromLastTransaction  float64 `json:"distance_from_last_transaction"` // Distance from the last transaction
	RatioToMedianPurchasePrice   float64 `json:"ratio_to_median_purchase_price"` // Ratio to median purchase price
	RepeatRetailer               bool    `json:"repeat_retailer"`                // Is it a repeat retailer?
	UsedChip                     bool    `json:"used_chip"`                      // Was a chip used?
	UsedPinNumber                bool    `json:"used_pin_number"`                // Was a PIN used?
	OnlineOrder                  bool    `json:"online_order"`                   // Was it an online order?
	Fraud                        bool    `json:"fraud"`                          // Is it a fraudulent transaction?
}
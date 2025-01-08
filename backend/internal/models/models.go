package models

import "time"

type Record struct {
	ID           uint      `gorm:"primaryKey" json:"-"` // Auto-increment ID
	AssetType    string    `json:"asset_type"`          // Crypto or ETF
	Timestamp    time.Time `json:"timestamp"`
	BidAskSpread float64   `json:"bid_ask_spread"` // Difference between ask and bid prices
	Volume       float64   `json:"volume"`         // Trading volume
	BidPrice     float64   `json:"bid_price"`      // High price (useful for trend analysis)
}

type LiquidityReport struct {
	AssetType                  string   `json:"asset_type"`
	TotalRecords               int      `json:"total_records"`
	HistoricalRecords          int      `json:"historical_records"`
	PredictionRecords          int      `json:"prediction_records"`
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
	ID                          uint    `gorm:"primaryKey" json:"id,omitempty"`
	DistanceFromHome            float64 `json:"distance_from_home"`
	DistanceFromLastTransaction float64 `json:"distance_from_last_transaction"`
	RatioToMedianPurchasePrice  float64 `json:"ratio_to_median_purchase_price"`
	RepeatRetailer              bool    `json:"repeat_retailer"`
	UsedChip                    bool    `json:"used_chip"`
	UsedPinNumber               bool    `json:"used_pin_number"`
	OnlineOrder                 bool    `json:"online_order"`
	Fraud                       bool    `json:"fraud"`
}

type TokenTransaction struct {
	BlockNumber string    `json:"blockNumber"` 
	TimeStamp   string    `json:"timeStamp"`   
	DateTime    time.Time `json:"dateTime"`           
	Hash        string    `json:"hash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       string    `json:"value"`
	TokenName   string    `json:"tokenName"`
	TokenSymbol string    `json:"tokenSymbol"`
}

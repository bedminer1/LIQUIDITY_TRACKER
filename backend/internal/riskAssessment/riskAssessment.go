package riskassessment

import (
	"fmt"

	"github.com/bedminer1/liquidity_tracker/internal/models"
)

func AssessLiquidity(currentRecords, predictions []models.Record, windowSize int) models.LiquidityReport {
	var report models.LiquidityReport
	if len(currentRecords) > 0 { report.AssetType = currentRecords[0].AssetType }

	allRecords := append(currentRecords, predictions...)

	// Sliding window for moving averages
	var volumeWindow []float64
	var spreadWindow []float64

	// Helper to calculate moving average
	movingAverage := func(data []float64) float64 {
		if len(data) == 0 {
			return 0
		}
		sum := 0.0
		for _, v := range data {
			sum += v
		}
		return sum / float64(len(data))
	}

	// Counters for risk levels
	currentHighRiskCount := 0
	currentModerateRiskCount := 0
	predictedHighRiskCount := 0
	predictedModerateRiskCount := 0

	var currentWarnings []string
	var predictedWarnings []string

	for idx, record := range allRecords {
		isPrediction := idx >= len(currentRecords)

		// Calculate severity
		volumeWindow = append(volumeWindow, record.Volume)
		spreadWindow = append(spreadWindow, record.BidAskSpread)

		if len(volumeWindow) > windowSize {
			volumeWindow = volumeWindow[1:]
		}
		if len(spreadWindow) > windowSize {
			spreadWindow = spreadWindow[1:]
		}

		volumeMA := movingAverage(volumeWindow)
		spreadMA := movingAverage(spreadWindow)

		isHighRisk := record.BidAskSpread > 3.0*spreadMA || record.Volume < 0.4*volumeMA
		isModerateRisk := record.BidAskSpread > 1.2*spreadMA || record.Volume < 0.7*volumeMA

		if isHighRisk {
			if isPrediction {
				predictedHighRiskCount++
				predictedWarnings = append(predictedWarnings, fmt.Sprintf("Predicted high risk for %s at %s: Spread=%.2f (MA=%.2f), Volume=%.0f (MA=%.0f)",
					record.AssetType, record.Timestamp, record.BidAskSpread, spreadMA, record.Volume, volumeMA))
			} else {
				currentHighRiskCount++
				currentWarnings = append(currentWarnings, fmt.Sprintf("Current high risk for %s at %s: Spread=%.2f (MA=%.2f), Volume=%.0f (MA=%.0f)",
					record.AssetType, record.Timestamp, record.BidAskSpread, spreadMA, record.Volume, volumeMA))
			}
		} else if isModerateRisk {
			if isPrediction {
				predictedModerateRiskCount++
			} else {
				currentModerateRiskCount++
			}
		}
	}

	report.TotalRecords = len(allRecords)
	report.HighRiskCount = currentHighRiskCount + predictedHighRiskCount
	report.ModerateRiskCount = currentModerateRiskCount + predictedModerateRiskCount

	// split into current n predicted records
	report.CurrentWarnings = currentWarnings
	report.PredictedWarnings = predictedWarnings
	report.CurrentHighRiskCount = currentHighRiskCount
	report.PredictedHighRiskCount = predictedHighRiskCount
	report.CurrentModerateRiskCount = currentModerateRiskCount
	report.PredictedModerateRiskCount = predictedModerateRiskCount

	return report
}

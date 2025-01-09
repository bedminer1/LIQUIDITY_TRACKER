package stats

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/bedminer1/liquidity_tracker/internal/models"
)

func GeneratePredictions(records []models.Record, intervals int, intervalLength int) []models.Record {
	if len(records) == 0 {
		return nil // No historical data to base predictions on
	}

	lastRecord := records[len(records)-1]
	lastTimestamp := lastRecord.Timestamp

	// Calculate trend and seasonality
	spreads := extractField(records, func(r models.Record) float64 { return r.BidAskSpread })
	movingAverage := calculateMovingAverage(spreads, 30)

	trendVolume := calculateTrend(records, func(r models.Record) float64 { return r.Volume })
	seasonalityVolume := calculateSeasonality(records, func(r models.Record) float64 { return r.Volume })

	// Calculate historical range for scaling noise
	minVolume, maxVolume := calculateRange(records, func(r models.Record) float64 { return r.Volume })
	volumeRange := maxVolume - minVolume

	var predictions []models.Record

	// Generate predictions
	for i := 1; i <= intervals; i++ {
		predictedTimestamp := lastTimestamp.Add(time.Duration(intervalLength*i) * time.Second)

		// Spread prediction
		predictedSpread := movingAverage * (0.9999 + rand.Float64()*0.0002)
		randomChoice := rand.Float64() // Random value between 0 and 1

		if randomChoice < 0.05 {
			// 10% chance of a spike (increase by 80%-200%)
			spikeFactor := 1.1 + rand.Float64()*1.6 // Random multiplier between 0.8 and 2.0
			predictedSpread *= spikeFactor
		}

		// Volume prediction
		volumeSeasonalIndex := i % len(seasonalityVolume)
		predictedVolume := lastRecord.Volume + trendVolume*float64(i) + seasonalityVolume[volumeSeasonalIndex]
		volumeVolatility := (rand.Float64() - 0.5) * volumeRange * 0.2 // ±5% noise
		predictedVolume = math.Max(predictedVolume+volumeVolatility, minVolume) // Ensure non-negative

		// Bid price prediction
		bidPriceVolatility := (rand.Float64() - 0.5) * lastRecord.BidPrice * 0.2 // ±2% noise
		predictedBidPrice := math.Max(lastRecord.BidPrice+bidPriceVolatility, 0)

		// Append prediction
		predictions = append(predictions, models.Record{
			AssetType:    lastRecord.AssetType,
			Timestamp:    predictedTimestamp,
			BidAskSpread: predictedSpread,
			Volume:       predictedVolume,
			BidPrice:     predictedBidPrice,
		})
	}

	return predictions
}

// Helper function to calculate a moving average
func calculateMovingAverage(data []float64, windowSize int) float64 {
	if len(data) < windowSize {
		windowSize = len(data)
	}

	sum := 0.0
	for i := len(data) - windowSize; i < len(data); i++ {
		sum += data[i]
	}

	return sum / float64(windowSize)
}

// Helper function to extract a field from records
func extractField(records []models.Record, selector func(models.Record) float64) []float64 {
	var result []float64
	for _, record := range records {
		result = append(result, selector(record))
	}
	return result
}

// Helper function to calculate the trend
func calculateTrend(records []models.Record, fieldSelector func(models.Record) float64) float64 {
	if len(records) < 2 {
		return 0
	}
	first := fieldSelector(records[0])
	last := fieldSelector(records[len(records)-1])
	return (last - first) / float64(len(records))
}

// Helper function to calculate seasonality
func calculateSeasonality(records []models.Record, fieldSelector func(models.Record) float64) []float64 {
	seasonLength := min(10, rand.IntN(int(math.Min(20, float64(len(records)))))) // Use up to 10 periods for seasonality
	seasonality := make([]float64, seasonLength)

	for i := 0; i < seasonLength; i++ {
		seasonality[i] = fieldSelector(records[i]) - calculateTrend(records, fieldSelector)*float64(i)
	}

	return seasonality
}

// Helper function to calculate the range of a field
func calculateRange(records []models.Record, fieldSelector func(models.Record) float64) (float64, float64) {
	minVal, maxVal := math.MaxFloat64, -math.MaxFloat64
	for _, r := range records {
		value := fieldSelector(r)
		if value < minVal {
			minVal = value
		}
		if value > maxVal {
			maxVal = value
		}
	}
	return minVal, maxVal
}
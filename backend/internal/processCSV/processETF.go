package processcsv

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bedminer1/liquidity_tracker/internal/models"
)

// Helper function to replace commas with periods for decimal parsing
func parseFloatWithComma(value string) (float64, error) {
	normalized := strings.ReplaceAll(value, ",", ".")
	return strconv.ParseFloat(normalized, 64)
}

func ParseEtfCsv(filePath string) ([]models.Record, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // ETF files use semicolons as delimiters
	reader.Read()      // Skip the header row

	var records []models.Record
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}

		fileName := filepath.Base(filePath)
		// Parse relevant fields
		date := line[0] + "T00:00:00Z"
		parsedDate, _ := time.Parse("02.01.2006T15:04:05Z", date)
		
		volume, err := parseFloatWithComma(line[5])
		if err != nil {
			volume = 0.0 // Default to 0 if invalid
		}
		bidAskSpreadPercent, err := parseFloatWithComma(line[6])
		if err != nil {
			bidAskSpreadPercent = 0.0 // Default to 0 if invalid
		}
		bidPrice, err := parseFloatWithComma(line[1])
		if err != nil {
			bidPrice = 0.0 // Default to 0 if invalid
		}

		// Convert bid-ask spread percentage to actual value
		bidAskSpread := bidAskSpreadPercent / 100.0

		records = append(records, models.Record{
			AssetType:    "ETF_" + fileName[:3],
			Timestamp:    parsedDate,
			BidAskSpread: bidAskSpread,
			Volume:       volume,
			BidPrice:     bidPrice,
		})
	}
	return records, nil
}

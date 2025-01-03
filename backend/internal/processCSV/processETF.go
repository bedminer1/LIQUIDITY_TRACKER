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
	fileName := filepath.Base(filePath)
	assetType := "ETF_" + fileName[:3]

	var records []models.Record
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}
		// Parse relevant fields
		date := line[0] + "T00:00:00Z"
		parsedDate, _ := time.Parse("02.01.2006T15:04:05Z", date)
		
		volume, err := parseFloatWithComma(line[5])
		if err != nil {
			volume = 0.0
		}
		bidAskSpread, err := parseFloatWithComma(line[6])
		if err != nil {
			continue // skip when it's NA
		}
		bidPrice, err := parseFloatWithComma(line[1])
		if err != nil {
			continue
		}

		records = append(records, models.Record{
			AssetType:    assetType,
			Timestamp:    parsedDate,
			BidAskSpread: bidAskSpread,
			Volume:       volume,
			BidPrice:     bidPrice,
		})
	}
	return records, nil
}

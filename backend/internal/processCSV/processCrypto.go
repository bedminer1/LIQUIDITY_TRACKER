package processcsv

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bedminer1/liquidity_tracker/internal/models"
)

func parseUnixTimestamp(unixTimestamp string) (time.Time, error) {
	// Split the string into seconds and fractional parts
	parts := strings.Split(unixTimestamp, ".")
	seconds, err := strconv.ParseInt(parts[0], 10, 64) // Parse seconds
	if err != nil {
		return time.Time{}, err
	}

	// Parse nanoseconds if fractional part exists
	nanoseconds := int64(0)
	if len(parts) > 1 {
		fractionalPart := parts[1]
		if len(fractionalPart) > 9 {
			fractionalPart = fractionalPart[:9] // Trim to nanoseconds precision
		}
		for len(fractionalPart) < 9 {
			fractionalPart += "0" // Pad with zeros to nanoseconds precision
		}
		nanoseconds, err = strconv.ParseInt(fractionalPart, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
	}

	// Convert to time.Time
	return time.Unix(seconds, nanoseconds).UTC(), nil
}

func ParseCryptoTxt(filePath string) ([]models.Record, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Regular expression to split by arbitrary whitespace
	re := regexp.MustCompile(`\s+`)
	assetName := strings.ToUpper(filepath.Base(filePath)[:3])

	var records []models.Record
	scanner := bufio.NewScanner(file)

	// Skip the header row
	if scanner.Scan() {
		// Header row is ignored
	}

	for scanner.Scan() {
		line := scanner.Text()
		fields := re.Split(line, -1) // Split by whitespace

		// Ensure the line has the correct number of fields
		if len(fields) < 11 {
			continue // Skip malformed lines
		}

		// Parse relevant fields
		unixTimestampStr := fields[0]
		timestamp, err := parseUnixTimestamp(unixTimestampStr)
		if err != nil {
			return nil, err
		}

		bidPrice, _ := strconv.ParseFloat(fields[1], 64)
		askPrice, _ := strconv.ParseFloat(fields[3], 64)
		volume, _ := strconv.ParseFloat(fields[8], 64)

		// Calculate bid-ask spread
		bidAskSpread := askPrice - bidPrice

		records = append(records, models.Record{
			AssetType:    "Crypto_" + assetName,
			Timestamp:    timestamp,
			BidAskSpread: bidAskSpread,
			Volume:       volume,
			BidPrice:         bidPrice,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

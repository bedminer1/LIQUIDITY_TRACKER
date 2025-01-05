package processcsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/bedminer1/liquidity_tracker/internal/models"
)

func ParseFraudCSV(filepath string) ([]models.TransactionRecord, error) {
	// Open the CSV file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Skip first row
	headers, err := reader.Read()
	fmt.Println("Headers:", headers) // Debugging log
	if err != nil {
		if err.Error() == "EOF" {
			return nil, fmt.Errorf("CSV file is empty or contains only headers")
		}
		return nil, err
	}

	// Parse the rows
	records := []models.TransactionRecord{}
	rowCount := 0
	for {
		// for logging
		rowCount++
		if rowCount % 100000 == 0 {
			fmt.Println("Parsing rows: ", rowCount, "/1000000")
		}

		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}

		// Convert row to TransactionRecord
		record, err := parseTransactionRow(row)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

// Helper function to parse a single CSV row
func parseTransactionRow(row []string) (models.TransactionRecord, error) {
	if len(row) < 8 {
		return models.TransactionRecord{}, errors.New("row has insufficient columns")
	}
	distanceFromHome, err := strconv.ParseFloat(row[0], 64)
	if err != nil {
		return models.TransactionRecord{}, err
	}
	distanceFromLastTransaction, err := strconv.ParseFloat(row[1], 64)
	if err != nil {
		return models.TransactionRecord{}, err
	}
	ratioToMedianPurchasePrice, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return models.TransactionRecord{}, err
	}
	repeatRetailer, err := parseBoolFromFloat(row[3])
	if err != nil {
		return models.TransactionRecord{}, err
	}
	usedChip, err := parseBoolFromFloat(row[4])
	if err != nil {
		return models.TransactionRecord{}, err
	}
	usedPinNumber, err := parseBoolFromFloat(row[5])
	if err != nil {
		return models.TransactionRecord{}, err
	}
	onlineOrder, err := parseBoolFromFloat(row[6])
	if err != nil {
		return models.TransactionRecord{}, err
	}
	fraud, err := parseBoolFromFloat(row[7])
	if err != nil {
		return models.TransactionRecord{}, err
	}
	return models.TransactionRecord{
		DistanceFromHome:            distanceFromHome,
		DistanceFromLastTransaction: distanceFromLastTransaction,
		RatioToMedianPurchasePrice:  ratioToMedianPurchasePrice,
		RepeatRetailer:              repeatRetailer,
		UsedChip:                    usedChip,
		UsedPinNumber:               usedPinNumber,
		OnlineOrder:                 onlineOrder,
		Fraud:                       fraud,
	}, nil
}

// Helper function to parse a float-based boolean (1.0 = true, 0.0 = false)
func parseBoolFromFloat(value string) (bool, error) {
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false, err
	}
	return floatVal == 1.0, nil
}

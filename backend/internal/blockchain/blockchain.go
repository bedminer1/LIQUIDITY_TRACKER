package blockchain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/bedminer1/liquidity_tracker/internal/models"
)

type APIResponse struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Result  []models.TokenTransaction `json:"result"`
}

func FetchTokenTransactions(contractAddress, walletAddress, apiKey string, start, end time.Time) ([]models.TokenTransaction, error) {
	// Construct the URL
	url := fmt.Sprintf(
		"https://api.etherscan.io/api?module=account&action=tokentx&contractaddress=%s&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s",
		contractAddress, walletAddress, apiKey,
	)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token transactions: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %v", err)
	}

	if apiResponse.Status != "1" {
		return nil, fmt.Errorf("API error: %s", apiResponse.Message)
	}

	filteredTransactions := []models.TokenTransaction{}
	for _, record := range apiResponse.Result {
		recordDateTime := unixToTime(record.TimeStamp)
		if recordDateTime.Before(start) {
			continue
		}
		if recordDateTime.After(end) {
			break
		}

		record.DateTime = recordDateTime
		filteredTransactions = append(filteredTransactions, record)
	}

	return filteredTransactions, nil
}

func unixToTime(unixTimestamp string) time.Time {
	timestamp, err := strconv.ParseInt(unixTimestamp, 10, 64)
	if err != nil {
		fmt.Printf("Error parsing timestamp: %v\n", err)
		return time.Time{}
	}
	return time.Unix(timestamp, 0).UTC()
}

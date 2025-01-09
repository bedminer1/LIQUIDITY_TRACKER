package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bedminer1/liquidity_tracker/internal/blockchain"
	"github.com/bedminer1/liquidity_tracker/internal/chatgpt"
	"github.com/bedminer1/liquidity_tracker/internal/models"
	riskassessment "github.com/bedminer1/liquidity_tracker/internal/riskAssessment"
	"github.com/bedminer1/liquidity_tracker/internal/stats"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func initHandler() *handler {
	db, err := gorm.Open(sqlite.Open("../../market_data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate()

	return &handler{DB: db}
}

func (h *handler) handleGetRecords(c echo.Context) error {
	asset, start, end, _, _, err := parseQueryParams(c)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": err,
		})
	}
	records, err := fetchRecordsFromDB(h.DB, asset, start, end)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": err,
		})
	}
	return c.JSON(200, echo.Map{
		"records": records,
	})
}


func (h *handler) handleGetBlockchainData(c echo.Context) error {
	contractAddress := c.QueryParam("contact_address")
	if contractAddress == "" {
		contractAddress = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48" // USDC
	}
	walletAddress := c.QueryParam("wallet_address")
	if walletAddress == "" {
		walletAddress = "0x28c6c06298d514db089934071355e5743bf21d60"
	}
	apiKey := c.QueryParam("api_key")
	if apiKey == "" {
		godotenv.Load("../../.env")
		apiKey = os.Getenv("ETHERSCAN_API_KEY")
	}
	start := c.QueryParam("start")
	if start == "" {
		start = "1974-01-01"
	}
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return fmt.Errorf("invalid 'start' date format, use YYYY-MM-DD")
	}
	end := c.QueryParam("end")
	if end == "" {
		end = "2030-01-01"
	}
	endTime, err := time.Parse("2006-01-02", end)
	if err != nil {
		return fmt.Errorf("invalid 'end' date format, use YYYY-MM-DD")
	}

	transactions, err := blockchain.FetchTokenTransactions(contractAddress, walletAddress, apiKey, startTime, endTime)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"transactions": transactions,
	})
}

func (h *handler) handleGetPredictions(c echo.Context) error {
	asset, start, end, intervalLength, intervals, err := parseQueryParams(c)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": err,
		})
	}
	records, err := fetchRecordsFromDB(h.DB, asset, start, end)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": err,
		})
	}
	predictions, err := getPredictionsFromAI(records, intervalLength, intervals)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": fmt.Sprintf("error interacting with microservice: %s", err.Error()),
		})
	}

	return c.JSON(200, echo.Map{
		"historicalData": records,
		"predictions":    predictions,
	})
}

func (h *handler) handleGetReport(c echo.Context) error {
	asset, start, end, intervalLength, intervals, err := parseQueryParams(c)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": err,
		})
	}
	records, err := fetchRecordsFromDB(h.DB, asset, start, end)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": err,
		})
	}
	predictions, err := getPredictionsFromAI(records, intervalLength, intervals)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": fmt.Sprintf("error interacting with microservice: %s", err.Error()),
		})
	}
	liquidityReport := riskassessment.AssessLiquidity(records, predictions, 8)

	return c.JSON(200, echo.Map{
		"report": liquidityReport,
	})
}

// HELPER FUNCTIONS

func parseQueryParams(c echo.Context) (string, time.Time, time.Time, int, int, error) {
	asset := c.QueryParam("asset")
	start := c.QueryParam("start")
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return "", time.Time{}, time.Time{}, 0, 0, fmt.Errorf("invalid 'start' date format, use YYYY-MM")
	}
	end := c.QueryParam("end")
	endTime, err := time.Parse("2006-01-02", end)
	if err != nil {
		return "", time.Time{}, time.Time{}, 0, 0, fmt.Errorf("invalid 'end' date format, use YYYY-MM")
	}
	intervalLength, _ := strconv.Atoi(c.QueryParam("time_interval_length"))
	intervals, _ := strconv.Atoi(c.QueryParam("time_intervals"))

	return asset, startTime, endTime, intervalLength, intervals, nil
}

func fetchRecordsFromDB(db *gorm.DB, asset string, start, end time.Time) ([]models.Record, error) {
	records := []models.Record{}
	err := db.Where("asset_type = ? AND timestamp BETWEEN ? AND ?", asset, start, end).Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching records from database: %v", err)
	}
	return records, nil
}


func getPredictionsFromAI(currentRecords []models.Record, intervalLength, intervals int) ([]models.Record, error) {
	// Convert records to JSON
	jsonData, err := json.Marshal(currentRecords)
	if err != nil {
		return nil, fmt.Errorf("error marshalling current records: %v", err)
	}

	// Send POST request to AI microservice
	url := fmt.Sprintf("http://localhost:5433/predict?time_interval_length=%d&time_intervals=%d", intervalLength, intervals)
	if intervalLength == 0 || intervals == 0 {
		url = "http://localhost:5433/predict"
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error sending request to AI microservice: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(resp)
		return nil, fmt.Errorf("AI microservice error: %s", string(body))
	}

	// Parse response
	var predictions []models.Record
	err = json.NewDecoder(resp.Body).Decode(&predictions)
	if err != nil {
		return nil, fmt.Errorf("error decoding response from AI microservice: %v", err)
	}

	return predictions, nil
}

func (h *handler) handleGetChatGPTRecommendation(c echo.Context) error {
	asset, start, end, intervalLength, intervals, err := parseQueryParams(c)
	intervals *= (intervalLength/86400)

	if err != nil {
		return c.JSON(400, echo.Map{
			"error": err,
		})
	}
	records, err := fetchRecordsFromDB(h.DB, asset, start, end)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": err,
		})
	}

	// predictions, err := getPredictionsFromAI(records, intervalLength, intervals)
	// if err != nil {
	// 	return c.JSON(400, echo.Map{
	// 		"error": fmt.Sprintf("error interacting with microservice: %s", err.Error()),
	// 	})
	// }

	predictions := stats.GeneratePredictions(records, intervals)
	liquidityReport := riskassessment.AssessLiquidity(records, predictions, 8)
	response, err := chatgpt.FetchGPTResponse(liquidityReport)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(201, echo.Map{
		"analysis": response.Choices[0].Message.Content,
		"report": liquidityReport,
		"historical_data": records,
		"predictions": predictions,
	})
}
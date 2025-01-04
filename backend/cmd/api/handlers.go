package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bedminer1/liquidity_tracker/internal/models"
	riskassessment "github.com/bedminer1/liquidity_tracker/internal/riskAssessment"
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
	asset := c.QueryParam("asset")
	start := c.QueryParam("start")
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return fmt.Errorf("invalid 'start' date format, use YYYY-MM")
	}
	end := c.QueryParam("end")
	endTime, err := time.Parse("2006-01-02", end)
	if err != nil {
		return fmt.Errorf("invalid 'end' date format, use YYYY-MM")
	}
	records := []models.Record{}
	h.DB.Where("asset_type = ? AND timestamp BETWEEN ? AND ?", asset, startTime, endTime).Find(&records)

	return c.JSON(200, echo.Map{
		"records": records,
	})
}

func (h *handler) handleGetPredictions(c echo.Context) error {
	// handle queries
	asset := c.QueryParam("asset")
	start := c.QueryParam("start")
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return fmt.Errorf("invalid 'start' date format, use YYYY-MM")
	}
	end := c.QueryParam("end")
	endTime, err := time.Parse("2006-01-02", end)
	if err != nil {
		return fmt.Errorf("invalid 'end' date format, use YYYY-MM")
	}

	// fetch "current" data
	currentRecords := []models.Record{}
	h.DB.Where("asset_type = ? AND timestamp BETWEEN ? AND ?", asset, startTime, endTime).Find(&currentRecords)

	// get predictions from microservice
	predictions, err := getPredictionsFromAI(currentRecords)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": fmt.Sprintf("error interacting with microservice: %s", err.Error()),
		})
	}

	return c.JSON(200, echo.Map{
		"historicalData": currentRecords,
		"predictions": predictions,
	})
}

func (h *handler) handleGetReport(c echo.Context) error {
	// handle queries
	asset := c.QueryParam("asset")
	start := c.QueryParam("start")
	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return fmt.Errorf("invalid 'start' date format, use YYYY-MM")
	}
	end := c.QueryParam("end")
	endTime, err := time.Parse("2006-01-02", end)
	if err != nil {
		return fmt.Errorf("invalid 'end' date format, use YYYY-MM")
	}

	// fetch "current" data
	currentRecords := []models.Record{}
	h.DB.Where("asset_type = ? AND timestamp BETWEEN ? AND ?", asset, startTime, endTime).Find(&currentRecords)

	// get predictions from microservice
	predictions, err := getPredictionsFromAI(currentRecords)
	if err != nil {
		return fmt.Errorf("error fetching predictions: %v", err)
	}

	// assess for liquidity shortfalls
	liquidityReport := riskassessment.AssessLiquidity(currentRecords, predictions, 8)

	return c.JSON(200, echo.Map{
		"report": liquidityReport,
	})
}

func getPredictionsFromAI(currentRecords []models.Record) ([]models.Record, error) {
	// Convert records to JSON
	jsonData, err := json.Marshal(currentRecords)
	if err != nil {
		return nil, fmt.Errorf("error marshalling current records: %v", err)
	}

	// Send POST request to AI microservice
	resp, err := http.Post("http://localhost:5433/predict", "application/json", bytes.NewBuffer(jsonData))
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

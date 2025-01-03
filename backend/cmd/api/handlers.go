package main

import (
	"fmt"
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
	var parseDate = func(dateStr string) time.Time {
		date, _ := time.Parse("01.02.2006", dateStr)
		return date
	}
	// mock data
	predictions := []models.Record{
		{AssetType: "ETF", Timestamp: parseDate("12.04.2013"), BidAskSpread: 0.0176, Volume: 806282, BidPrice: 120.6972},
		{AssetType: "ETF", Timestamp: parseDate("15.04.2013"), BidAskSpread: 0.0504, Volume: 462017, BidPrice: 120.51},
		{AssetType: "ETF", Timestamp: parseDate("16.04.2013"), BidAskSpread: 0.0506, Volume: 1375894, BidPrice: 120.35},
		{AssetType: "ETF", Timestamp: parseDate("17.04.2013"), BidAskSpread: 0.0227, Volume: 1258238, BidPrice: 119.85},
		{AssetType: "ETF", Timestamp: parseDate("18.04.2013"), BidAskSpread: 0.0356, Volume: 771955, BidPrice: 119.84},
		{AssetType: "ETF", Timestamp: parseDate("19.04.2013"), BidAskSpread: 0.0294, Volume: 353449, BidPrice: 119.85},
		{AssetType: "ETF", Timestamp: parseDate("22.04.2013"), BidAskSpread: 0.0898, Volume: 472687, BidPrice: 120.19},
		{AssetType: "ETF", Timestamp: parseDate("23.04.2013"), BidAskSpread: 0.0234, Volume: 700594, BidPrice: 120.48},
		{AssetType: "ETF", Timestamp: parseDate("24.04.2013"), BidAskSpread: 0.0201, Volume: 696207, BidPrice: 120.65},
	}

	// assess for liquidity shortfalls
	liquidityReport := riskassessment.AssessLiquidity(currentRecords, predictions, 8)

	return c.JSON(200, echo.Map{
		"report": liquidityReport,
	})
}

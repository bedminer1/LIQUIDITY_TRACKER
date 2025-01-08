package main

import (
	"github.com/labstack/echo/v4"
)

const version = "1.0.0"

func main() {
	e := echo.New()
	h := initHandler()
	
	e.GET("/healthcheck", h.handleHealthCheck)
	e.GET("/records", h.handleGetRecords)
	e.GET("/blockchain_records", h.handleGetBlockchainData)
	e.GET("/predictions", h.handleGetPredictions)
	e.GET("/report", h.handleGetReport)
	e.GET("/recommendations", h.handleGetChatGPTRecommendation)

	e.Logger.Fatal(e.Start(":4000"))
}
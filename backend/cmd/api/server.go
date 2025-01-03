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
	e.GET("/report", h.handleGetReport)

	e.Logger.Fatal(e.Start(":4000"))
}
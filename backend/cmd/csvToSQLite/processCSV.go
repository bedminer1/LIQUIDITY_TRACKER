package main

import (
	"log"

	"github.com/bedminer1/liquidity_tracker/internal/models"
	processcsv "github.com/bedminer1/liquidity_tracker/internal/processCSV"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("../../market_data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the Record struct
	db.AutoMigrate(&models.Record{})
	return db
}

func insertRecords(db *gorm.DB, records []models.Record) error {
	for _, record := range records {
		if err := db.Create(&record).Error; err != nil {
			return err
		}
	}
	return nil
}

func main() {
	db := initDB()
	cryptoRecords, err := processcsv.ParseCryptoTxt("../../data/crypto_data/btcusd.txt")
	if err != nil {
		log.Fatal("Error parsing crypto CSV:", err)
	}

	etfRecords, err := processcsv.ParseEtfCsv("../../data/etf_data/EMB_data.csv")
	if err != nil {
		log.Fatal("Error parsing ETF CSV:", err)
	}

	allRecords := append(cryptoRecords, etfRecords...)
	if err := insertRecords(db, allRecords); err != nil {
		log.Fatal("Error inserting records into database:", err)
	}

	log.Println("Data successfully inserted into the database.")
}

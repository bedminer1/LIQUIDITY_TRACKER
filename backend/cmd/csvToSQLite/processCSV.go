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
	records := []models.Record{}

	etfFileNames := []string{
		"EMB_data.csv",
		"HYG_data.csv",
		"LQD_data.csv",
		"TLT_data.csv",
	}
	cryptoFileNames := []string{
		"btcusd.txt",
		"ethusd.txt",
		"xrpusd.txt",
	}
	etfPrefix := "../../data/etf_data/"
	cryptoPrefix := "../../data/crypto_data/"

	for _, cryptoFileName := range cryptoFileNames {
		filePath := cryptoPrefix + cryptoFileName
		cryptoRecords, err := processcsv.ParseCryptoTxt(filePath)
		if err != nil {
			log.Fatal("Error parsing crypto CSV:", err)
		}
		records = append(records, cryptoRecords...)
	}

	for _, etfFileName := range etfFileNames {
		filePath := etfPrefix + etfFileName
		etfRecords, err := processcsv.ParseEtfCsv(filePath)
		if err != nil {
			log.Fatal("Error parsing ETF CSV:", err)
		}
		records = append(records, etfRecords...)
	}

	if err := insertRecords(db, records); err != nil {
		log.Fatal("Error inserting records into database:", err)
	}

	log.Println("Data successfully inserted into the database.")
}

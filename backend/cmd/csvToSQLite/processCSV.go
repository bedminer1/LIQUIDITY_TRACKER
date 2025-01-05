package main

import (
	"fmt"
	"log"

	"github.com/bedminer1/liquidity_tracker/internal/models"
	processcsv "github.com/bedminer1/liquidity_tracker/internal/processCSV"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("../../market_data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable logging
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the Record and TransactionRecord structs
	db.AutoMigrate(&models.Record{}, &models.TransactionRecord{})
	return db
}

// Generalized function to insert records of different types
func insertRecords[T any](db *gorm.DB, records []T) error {
	for i, record := range records {
		if i%50000 == 0 {
			fmt.Println("Inserting records: ", i+1, "/", len(records))
		}
		if err := db.Create(&record).Error; err != nil {
			return err
		}
	}
	return nil
}

func main() {
	db := initDB()

	// Insert Market Records
	// marketRecords := getMarketRecords() // Replace with your method to get market records
	// if err := insertRecords(db, marketRecords); err != nil {
	// 	log.Fatal("Error inserting market records into database:", err)
	// }
	// log.Println("Market data successfully inserted into the database.")

	// Insert Transaction Records
	transactionRecords := getTransactionRecords() // Replace with your method to get transaction records
	if err := insertRecords(db, transactionRecords); err != nil {
		log.Fatal("Error inserting transaction records into database:", err)
	}
	log.Println("Transaction data successfully inserted into the database.")
}

func getTransactionRecords() []models.TransactionRecord {
	records, err := processcsv.ParseFraudCSV("../../data/fraud_data/card_transdata.csv")
	if err != nil {
		log.Fatal("Error parsing fraud records:", err)
	}
	return records
}

package main

import (
	"log"
	"sync"

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

	var wg sync.WaitGroup
	recordChan := make(chan []models.Record, len(etfFileNames)+len(cryptoFileNames))
	errChan := make(chan error, len(etfFileNames)+len(cryptoFileNames))

	for _, cryptoFileName := range cryptoFileNames {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()

			filePath := cryptoPrefix + fileName
			cryptoRecords, err := processcsv.ParseCryptoTxt(filePath)
			if err != nil {
				errChan <- err
				return
			}
			recordChan <- cryptoRecords
		}(cryptoFileName)
	}

	for _, etfFileName := range etfFileNames {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()

			filePath := etfPrefix + fileName
			etfRecords, err := processcsv.ParseEtfCsv(filePath)
			if err != nil {
				errChan <- err
				return
			}
			recordChan <- etfRecords
		}(etfFileName)
	}

	go func() {
		wg.Wait()
		close(recordChan)
		close(errChan)
	}()

	for {
		select {
		case rec, ok := <-recordChan:
			if ok {
				records = append(records, rec...)
			}
		case err, ok := <-errChan:
			if ok {
				log.Fatal("Error processing file:", err)
			}
		default:
			// Exit loop when all channels are closed
			if len(recordChan) == 0 && len(errChan) == 0 {
				goto DONE
			}
		}
	}

DONE:
	if err := insertRecords(db, records); err != nil {
		log.Fatal("Error inserting records into database:", err)
	}

	log.Println("Data successfully inserted into the database.")
}

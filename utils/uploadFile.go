package utils

import (
	"Mymodule/models"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func UploadFileDB(ctx *gin.Context) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting DB object from GORM: %v", err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	file, _, err := ctx.Request.FormFile("User_data")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file", "details": err.Error()})
		return
	}

	if file == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file received"})
		return
	}

	// Limit file size to 30GB
	const maxFileSize = 30 << 30
	fileSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking file size", "details": err.Error()})
		return
	}
	if fileSize > maxFileSize {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the limit."})
		return
	}

	// Reset file pointer
	file.Seek(0, io.SeekStart)

	chunkSize := int64(10 * 1024 * 1024) // 10 MB chunks
	batchSize := 5000                    // Batch size for database inserts

	if err := uploadFileToDB(file, db, chunkSize, batchSize, fileSize); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing file", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded and processed successfully."})
}
func uploadFileToDB(file multipart.File, db *gorm.DB, chunkSize int64, batchSize int, fileSize int64) error {
	buffer := make([]byte, chunkSize) // Buffer to hold large chunks of the CSV file
	totalRead := int64(0)
	batch := make([]models.UserDetails, 0, batchSize)
	recordCount := 0

	for {
		// Read a large chunk of data into the buffer
		n, err := file.Read(buffer)
		if n == 0 {
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error reading file: %v", err)
			}
		}

		totalRead += int64(n)
		records := strings.Split(string(buffer[:n]), "\n") // Split buffer into lines

		for _, record := range records {
			if record == "" {
				continue // Skip empty lines
			}

			fields := strings.Split(record, ",") // Split by commas to get CSV fields
			if len(fields) != 11 {
				log.Printf("Skipping record with incorrect number of fields: %v", fields)
				continue
			}

			// Map fields to the UserDetails struct
			idStr := fields[0]
			if idStr == "" {
				idStr = "11111111"
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Printf("Invalid ID: %v", err)
				continue
			}

			batch = append(batch, models.UserDetails{
				ID:         id,
				FirstName:  fields[1],
				LastName:   fields[2],
				Email:      fields[3],
				Age:        fields[4],
				Gender:     fields[5],
				Department: fields[6],
				Company:    fields[7],
				Salary:     fields[8],
				DateJoined: fields[9],
				IsActive:   fields[10],
			})
			recordCount++

			// If the batch is full, insert into the database
			if len(batch) >= batchSize {
				if err := db.CreateInBatches(batch, batchSize).Error; err != nil {
					return fmt.Errorf("error inserting batch into DB: %v", err)
				}
				runtime.GC() // Request garbage collection
				log.Printf("Inserted %d records into the database...", len(batch))
				batch = batch[:0] // Clear the batch for the next set of records
			}
		}

		// Handle the case where we've reached the end of the buffer but there may be leftover data
		// If the buffer isn't perfectly aligned to full records, we need to process the leftovers
		if n < len(buffer) {
			leftover := buffer[n:]
			if len(leftover) > 0 {
				// Process any leftover data that was not yet split into full records
				leftoverRecords := strings.Split(string(leftover), "\n")
				for _, leftoverRecord := range leftoverRecords {
					if leftoverRecord == "" {
						continue
					}
					fields := strings.Split(leftoverRecord, ",")
					if len(fields) != 11 {
						log.Printf("Skipping leftover record with incorrect number of fields: %v", fields)
						continue
					}

					idStr := fields[0]
					if idStr == "" {
						idStr = "11111111"
					}
					id, err := strconv.ParseInt(idStr, 10, 64)
					if err != nil {
						log.Printf("Invalid ID in leftover data: %v", err)
						continue
					}

					batch = append(batch, models.UserDetails{
						ID:         id,
						FirstName:  fields[1],
						LastName:   fields[2],
						Email:      fields[3],
						Age:        fields[4],
						Gender:     fields[5],
						Department: fields[6],
						Company:    fields[7],
						Salary:     fields[8],
						DateJoined: fields[9],
						IsActive:   fields[10],
					})
					recordCount++
				}
			}
		}

		// If the buffer is processed and not aligned perfectly, the loop will continue
		if totalRead >= fileSize {
			break
		}
	}

	// Insert any remaining records after the last chunk is processed
	if len(batch) > 0 {
		if err := db.CreateInBatches(batch, batchSize).Error; err != nil {
			return fmt.Errorf("error inserting final batch into DB: %v", err)
		}
		log.Printf("Inserted %d remaining records into the database...", len(batch))
	}

	log.Printf("Successfully processed %d records", recordCount)
	return nil
}

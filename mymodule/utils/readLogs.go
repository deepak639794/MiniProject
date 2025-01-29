package utils

import (
	"Mymodule/mymodule/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var DefaultLogReader LogReader = &logReaderImpl{}

// logReaderImpl is the concrete implementation of LogReader
type logReaderImpl struct{}

func NewLogReader() *logReaderImpl {
	return &logReaderImpl{}
}

func (l *logReaderImpl) ReadLogsFromFile() ([]models.ApiLog, error) {
	fmt.Println("inside the function ")
	logFile, err := os.OpenFile("mymodule/logs/api_logs.json", os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
		return nil, err
	}
	defer logFile.Close()

	// Slice to store the logs
	var logs []models.ApiLog
	fmt.Println("here it is ")
	decoder := json.NewDecoder(logFile)

	if err := decoder.Decode(&logs); err != nil {
		if err.Error() == "EOF" {
			return logs, nil
		}
		log.Printf("Error decoding log entries: %v", err)
		return nil, err
	}
	fmt.Println("at the end")

	// Return the successfully decoded logs
	return logs, nil
}

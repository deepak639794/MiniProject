package utils

import (
	"Mymodule/models"
	"encoding/json"
	"log"
	"os"
)

// Fetch logs from the file

// Read the file content and decode it as an array of logs
func ReadLogsFromFile() ([]models.ApiLog, error) {
	logFile, err := os.OpenFile("./logs/api_logs.json", os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
		return nil, err
	}
	defer logFile.Close()

	// Slice to store the logs
	var logs []models.ApiLog

	// Create a JSON decoder and decode the array of log entries from the file
	decoder := json.NewDecoder(logFile)

	// Try decoding the entire JSON file as an array of logs
	if err := decoder.Decode(&logs); err != nil {
		// Handle EOF gracefully (empty file)
		if err.Error() == "EOF" {
			// File is empty, return empty slice
			return logs, nil
		}
		// If any other error, print it and return
		log.Printf("Error decoding log entries: %v", err)
		return nil, err
	}

	// Return the successfully decoded logs
	return logs, nil
}

// func ReadLogsFromFile() ([]models.ApiLog, error) {
// 	var logs []models.ApiLog

// 	// Open the log file
// 	fmt.Println("inside the function ")
// 	logFile, err := os.Open("api_logs.json")
// 	if err != nil {
// 		fmt.Println("some issues in opening the file ")
// 		return nil, err
// 	}
// 	defer logFile.Close()

// 	// Read the file line by line (each line represents a log entry)
// 	decoder := json.NewDecoder(logFile)
// 	for {
// 		fmt.Println("inside the for loop")
// 		var logEntry models.ApiLog
// 		if err := decoder.Decode(&logEntry); err != nil {
// 			fmt.Println("in the if condition")
// 			if err.Error() == "EOF" {
// 				fmt.Println("in the if condition")
// 				break // End of file
// 			}
// 			fmt.Println("there is an error ", err.Error())
// 			return nil, err
// 		}
// 		fmt.Println("outside the loop")
// 		fmt.Println(logs)
// 		logs = append(logs, logEntry)
// 	}

// 	return logs, nil
// }

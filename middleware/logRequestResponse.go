// middleware/logging.go
package middleware

import (
	"Mymodule/models"
	"bytes"
	"encoding/json"

	"io"
	"log"
	"os"
	"time" // Import your models package

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LogRequestResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a unique request ID for this request
		requestID := uuid.New().String()

		// Capture request details
		start := time.Now()
		method := c.Request.Method
		url := c.Request.URL.Path

		// Read request body (if any)
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Reset the body so it can be read later
		}

		// Capture the response using a custom writer
		respBody := &bytes.Buffer{}
		c.Writer = &CustomResponseWriter{
			ResponseWriter: c.Writer,
			Body:           respBody,
		}

		// Process the request
		c.Next()

		// After request processing: capture response status and body
		responseCode := c.Writer.Status()
		responseBody := respBody.String()

		// Create the log entry
		logEntry := models.ApiLog{
			RequestMethod: method,
			RequestURL:    url,
			RequestBody:   requestBody,
			ResponseCode:  responseCode,
			ResponseBody:  responseBody,
			CreatedAt:     start,
			RequestID:     requestID,
		}

		// Log the entry to a file in JSON format
		writeLogToFile(logEntry)
		wrapLogsInArray()

		// Optionally, log the details to console for debugging purposes
		log.Printf("RequestID: %s, Method: %s, URL: %s, ResponseCode: %d, Duration: %v",
			requestID, method, url, responseCode, time.Since(start))
	}
}

// CustomResponseWriter is used to capture response body
type CustomResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(p []byte) (int, error) {
	// Capture the response body
	w.Body.Write(p)
	return w.ResponseWriter.Write(p)

}
func writeLogToFile(logEntry models.ApiLog) {
	// Open the log file in read-write mode, or create it if it doesn't exist
	logFile, err := os.OpenFile("./logs/api_logs.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Read the existing logs from the file as an array
	var logs []models.ApiLog
	decoder := json.NewDecoder(logFile)

	// If the file is not empty, decode its contents into the logs array
	if err := decoder.Decode(&logs); err != nil && err.Error() != "EOF" {
		log.Printf("Error decoding log entries: %v", err)
		return
	}

	// Append the new log entry to the logs slice
	logs = append(logs, logEntry)

	// Truncate the file and seek to the beginning for rewriting
	if err := logFile.Truncate(0); err != nil {
		log.Printf("Error truncating log file: %v", err)
		return
	}

	if _, err := logFile.Seek(0, 0); err != nil {
		log.Printf("Error seeking to the beginning of the file: %v", err)
		return
	}

	// Encode and write the updated logs back to the file as a JSON array
	enc := json.NewEncoder(logFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(logs); err != nil {
		log.Printf("Error encoding logs back to file: %v", err)
		return
	}

	log.Printf("Log entry successfully appended to the file.")
}

// writeLogToFile writes the log entry to a file in JSON format
// func writeLogToFile(logEntry models.ApiLog) {
// 	// Open log file (append mode)
// 	logFile, err := os.OpenFile("api_logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Fatalf("Error opening log file: %v", err)
// 	}
// 	defer logFile.Close()

// 	// Convert log entry to JSON
// 	logData, err := json.Marshal(logEntry)
// 	if err != nil {
// 		log.Printf("Error marshalling log entry: %v", err)
// 		return
// 	}

//		// Write the log entry to the file followed by a newline
//		logFile.Write(append(logData, '\n'))
//	}
func wrapLogsInArray() {
	// Read the existing logs from the file
	logFile, err := os.OpenFile("./logs/api_logs.json", os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Read the file content
	var logs []models.ApiLog
	decoder := json.NewDecoder(logFile)
	for {
		var logEntry models.ApiLog
		if err := decoder.Decode(&logEntry); err != nil {
			if err.Error() == "EOF" {
				break // End of file
			}
			log.Printf("Error decoding log entry: %v", err)
			return
		}
		logs = append(logs, logEntry)
	}

	// Rewrite the logs to the file in JSON array format
	logFile.Truncate(0) // Clear the file
	logFile.Seek(0, 0)  // Reset the file pointer

	// Wrap logs in a JSON array
	logsArray := []models.ApiLog(logs)
	enc := json.NewEncoder(logFile)
	enc.SetIndent("", "  ")
	enc.Encode(logsArray)
}

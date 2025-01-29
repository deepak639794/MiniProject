// middleware/logging.go
package middleware

import (
	"Mymodule/mymodule/models"
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
		// responseBody := respBody.String()

		// Create the log entry
		logEntry := models.ApiLog{
			RequestMethod: method,
			RequestURL:    url,
			RequestBody:   requestBody,
			ResponseCode:  responseCode,
			ResponseBody:  "Success",
			CreatedAt:     start.String(),
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
	logFile, err := os.OpenFile("C:/Users/deepak.ag/Desktop/MiniProject1/mymodule/logs/api_logs.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening log file11: %v", err)
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

func wrapLogsInArray() {
	// Read the existing logs from the file
	logFile, err := os.OpenFile("C:/Users/deepak.ag/Desktop/MiniProject1/mymodule/logs/api_logs.json", os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Error opening log file1: %v", err)
	}
	defer logFile.Close()

	// Read the file content into a slice of ApiLog
	var logs []models.ApiLog
	decoder := json.NewDecoder(logFile)

	// Decode the root JSON array into the logs slice
	if err := decoder.Decode(&logs); err != nil {
		if err.Error() != "EOF" {
			log.Printf("Error decoding log entries: %v", err)
			return
		}
	}

	// Clear the file before writing the new content
	logFile.Truncate(0) // Clear the file
	logFile.Seek(0, 0)  // Reset the file pointer

	// Now rewrite the logs to the file in JSON array format
	enc := json.NewEncoder(logFile)
	enc.SetIndent("", "  ")

	// Write the logs as an array in the JSON file
	if err := enc.Encode(logs); err != nil {
		log.Printf("Error encoding logs back to file: %v", err)
		return
	}
}

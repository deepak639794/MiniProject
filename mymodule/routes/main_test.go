package routes

import (
	"Mymodule/mymodule/mocks"
	"Mymodule/mymodule/models"

	// mocks "Mymodule/mymodule/mock_utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHelloEndpoint(t *testing.T) {
	// Initialize Gin router
	r := gin.Default()

	SetupRoutes(r, nil, nil)

	req, _ := http.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check if the response status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the response body contains the expected message
	assert.JSONEq(t, `{"message": "hello sir"}`, w.Body.String())
}

// func TestFetchUsersEndpoint(t *testing.T) {
// 	r := gin.Default()

// 	routes.SetupRoutes(r, nil)

// 	req, _ := http.NewRequest("GET", "/hello", nil)
// 	w := httptest.NewRecorder()

// 	r.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Check if the response body matches the expected output
// 	expectedResponse := `{"message": "hello sir"}`
// 	assert.JSONEq(t, expectedResponse, w.Body.String())
// }

type MockUtils struct {
	mock.Mock
}

// Mock the ReadLogsFromFile function
func (m *MockUtils) ReadLogsFromFile() ([]models.ApiLog, error) {
	args := m.Called()
	return args.Get(0).([]models.ApiLog), args.Error(1)
}

func TestGetLogsEndpoint(t *testing.T) {
	// Mock the `ReadLogsFromFile` behavior
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock instance of LogReader
	mockReader := mocks.NewMockLogReader(ctrl)

	mockLogs := []models.ApiLog{
		{
			RequestMethod: "GET",
			RequestURL:    "/logs",
			RequestBody:   "",
			ResponseCode:  200,
			ResponseBody:  "Success",
			CreatedAt:     "2025-01-21T08:20:19.287970282Z",
			RequestID:     "ccfaf098-606f-4267-9977-7c2fab494bf3"},
	}

	mockReader.EXPECT().ReadLogsFromFile().Return(mockLogs, nil)
	router := gin.Default()
	SetupRoutes(router, mockReader, nil)

	// Create a test request
	req, _ := http.NewRequest("GET", "/logs", nil)
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Validate the response
	require.Equal(t, 200, w.Code)

	var response map[string][]models.ApiLog
	err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&response)
	require.NoError(t, err)

	logs := response["logs"]
	lastLog := logs[len(logs)-1]

	require.Equal(t, mockLogs[0].ResponseBody, lastLog.ResponseBody)

}

// models/log.go
package models

type ApiLog struct {
	RequestMethod string `json:"request_method"`
	RequestURL    string `json:"request_url"`
	RequestBody   string `json:"request_body"`
	ResponseCode  int    `json:"response_code"`
	ResponseBody  string `json:"response_body"`
	CreatedAt     string `json:"created_at"`
	RequestID     string `json:"request_id"`
}

package models

import "time"

type StatPaginatedResponse struct {
	Data  []*StatResponse `json:"data"`
	Total int64           `json:"total"`
}

type StatRequest struct {
	Limit     int64     `json:"limit"`
	Page      int64     `json:"page"`
	DateStart time.Time `json:"date-start"`
	DateEnd   time.Time `json:"date-end"`
}

type StatResponse struct {
	Ip             string    `json:"ip"`
	Bot            bool      `json:"bot"`
	Mobile         bool      `json:"mobile"`
	Browser        string    `json:"browser"`
	Os             string    `json:"os"`
	Platform       string    `json:"platform"`
	Referer        string    `json:"referer"`
	BrowserVersion string    `json:"browser_version"`
	Created        time.Time `json:"created"`
}

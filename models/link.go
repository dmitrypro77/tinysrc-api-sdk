package models

import "time"

type LinkRequest struct {
	Url            string `json:"url"`
	AuthRequired   int    `json:"auth_required"`
	Password       string `json:"password,omitempty"`
	ExpirationTime string `json:"expiration_time,omitempty"`
}

type ListUrlsRequest struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Query string `json:"query,omitempty"`
}

type LinkActivationRequest struct {
	Active bool `json:"active"`
}

type LinkResponse struct {
	Url          string `json:"url"`
	StatUrl      string `json:"stat_url,omitempty"`
	StatPassword string `json:"stat_password,omitempty"`
	Password     string `json:"password,omitempty"`
	AuthRequired int    `json:"auth_required"`
}

type LinkUserResponse struct {
	Url            string     `json:"url"`
	Hash           string     `json:"hash"`
	AuthRequired   int        `json:"auth_required"`
	Password       string     `json:"password,omitempty"`
	StatPassword   string     `json:"stat_password"`
	QRCode         string     `json:"qr_code"`
	Active         int        `json:"active"`
	Clicks         int64      `json:"clicks"`
	Bots           int64      `json:"bots"`
	StatUrl        string     `json:"stat_url"`
	Created        *time.Time `json:"created"`
	ExpirationTime *time.Time `json:"expiration_time,omitempty"`
}

type PaginatedLinkUserResponse struct {
	Data  []*LinkUserResponse `json:"data"`
	Total int64               `json:"total"`
}

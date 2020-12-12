package models

type CurrentUserResponse struct {
	Username string `json:"username"`
	ApiKey   string `json:"api_key"`
	Active   int    `json:"active"`
	Plan     int    `json:"plan"`
	Banned   int    `json:"banned"`
	Email    string `json:"email"`
}

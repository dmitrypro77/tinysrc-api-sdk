package models

type ErrorResponse struct {
	Validations map[string][]string `json:"validations"`
	Errors      []string            `json:"errors"`
	Status      int                 `json:"status"`
}

package models

type Exception struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Err string
}
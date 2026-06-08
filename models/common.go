package models

import "errors"

var (
	ErrCountryNameRequired = errors.New("country_name is required")
	ErrStatusRequired      = errors.New("status is required")
	ErrInvalidStatus       = errors.New("status must be 'Planned' or 'Visited'")
	ErrNotFound            = errors.New("record not found")
	ErrUnauthorized        = errors.New("unauthorized")
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type APIError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

type DashboardSummary struct {
	Total   int `json:"total"`
	Planned int `json:"planned"`
	Visited int `json:"visited"`
}

// NewSuccessResponse creates a standardized success API response
func NewSuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
}

// NewErrorResponse creates a standardized error API response
func NewErrorResponse(err string, code int) APIError {
	return APIError{
		Success: false,
		Error:   err,
		Code:    code,
	}
}

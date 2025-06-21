package utils

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		expected string
	}{
		{
			name: "error with underlying error",
			appError: &AppError{
				Code:    http.StatusBadRequest,
				Message: "Invalid input",
				Err:     errors.New("field is required"),
			},
			expected: "Invalid input: field is required",
		},
		{
			name: "error without underlying error",
			appError: &AppError{
				Code:    http.StatusNotFound,
				Message: "Resource not found",
				Err:     nil,
			},
			expected: "Resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.appError.Error(); got != tt.expected {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewAppError(t *testing.T) {
	err := errors.New("underlying error")
	appErr := NewAppError(http.StatusInternalServerError, "Something went wrong", err)

	if appErr.Code != http.StatusInternalServerError {
		t.Errorf("Expected code %d, got %d", http.StatusInternalServerError, appErr.Code)
	}

	if appErr.Message != "Something went wrong" {
		t.Errorf("Expected message 'Something went wrong', got '%s'", appErr.Message)
	}

	if appErr.Err != err {
		t.Errorf("Expected underlying error to be preserved")
	}
}

func TestErrorConstructors(t *testing.T) {
	err := errors.New("test error")

	tests := []struct {
		name         string
		constructor  func(string, error) *AppError
		expectedCode int
	}{
		{"BadRequestError", BadRequestError, http.StatusBadRequest},
		{"UnauthorizedError", UnauthorizedError, http.StatusUnauthorized},
		{"NotFoundError", NotFoundError, http.StatusNotFound},
		{"InternalServerError", InternalServerError, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := tt.constructor("test message", err)
			if appErr.Code != tt.expectedCode {
				t.Errorf("Expected code %d, got %d", tt.expectedCode, appErr.Code)
			}
			if appErr.Message != "test message" {
				t.Errorf("Expected message 'test message', got '%s'", appErr.Message)
			}
			if appErr.Err != err {
				t.Errorf("Expected underlying error to be preserved")
			}
		})
	}
}
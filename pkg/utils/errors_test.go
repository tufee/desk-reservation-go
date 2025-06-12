package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHTTPError(t *testing.T) {
	tests := []struct {
		name           string
		err           error
		expectedCode   int
		expectedBody   map[string]interface{}
		expectedStatus string
	}{
		{
			name:         "bad request error",
			err:         NewBadRequestError("invalid input"),
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "invalid input",
			},
		},
		{
			name:         "internal server error with wrapped error",
			err:         NewInternalServerError("processing failed", errors.New("database error")),
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"message": "processing failed: database error",
			},
		},
		{
			name:         "default error case",
			err:         errors.New("unknown error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"message": "Internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			HandleHTTPError(rr, tt.err)

			// Check status code
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}

			// Check Content-Type header
			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want %v",
					contentType, "application/json")
			}

			// Check response body
			var response map[string]interface{}
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Fatalf("Unable to decode response body: %v", err)
			}

			// Compare response with expected
			for key, expectedValue := range tt.expectedBody {
				if actualValue, exists := response[key]; !exists {
					t.Errorf("Expected key %s in response, but it was missing", key)
				} else if actualValue != expectedValue {
					t.Errorf("handler returned unexpected body: got %v want %v",
						actualValue, expectedValue)
				}
			}
		})
	}
}
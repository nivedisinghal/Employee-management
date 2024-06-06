package handlers

import (
	"ems/store"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestGetEmployeeHandler(t *testing.T) {
	// Initialize some employees for testing
	store.CreateEmployee("John Doe", "Developer", 60000.0)
	store.CreateEmployee("Alice Smith", "Manager", 80000.0)

	tests := []struct {
		name         string
		id           int
		expectedCode int
		expectedBody string
	}{
		// Valid case: Employee exists
		{
			name:         "Valid employee",
			id:           1,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"name":"John Doe","position":"Developer","salary":60000}`,
		},
		// Employee does not exist
		{
			name:         "Non-existent employee",
			id:           90,
			expectedCode: http.StatusNotFound,
			expectedBody: "Employee not found",
		},
		// Invalid ID: less than 1
		{
			name:         "Invalid ID",
			id:           0,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid employee ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/employees/"+strconv.Itoa(tt.id), nil)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(GetEmployeeHandler)

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					recorder.Code, tt.expectedCode)
			}

			// Check response body for error cases
			if recorder.Code != http.StatusOK {
				body := strings.TrimSpace(recorder.Body.String())
				if body != tt.expectedBody {
					t.Errorf("handler returned unexpected body: got %v want %v",
						body, tt.expectedBody)
				}
			}
		})
	}
}

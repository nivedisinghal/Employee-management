package handlers

import (
	"ems/models"
	"ems/store"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestListEmployeesHandler(t *testing.T) {
	// Initialize some employees for testing
	store.CreateEmployee("John Doe", "Developer", 60000.0)
	store.CreateEmployee("Alice Smith", "Manager", 80000.0)
	store.CreateEmployee("Bob Johnson", "Designer", 70000.0)
	store.CreateEmployee("Eve Williams", "Tester", 65000.0)

	tests := []struct {
		name              string
		page              int
		size              int
		expectedCode      int
		expectedBodyCount int
		expectedBody      string
	}{
		// Valid case: List first page with 2 employees per page
		{
			name:              "List first page with 2 employees per page",
			page:              1,
			size:              2,
			expectedCode:      http.StatusOK,
			expectedBodyCount: 2,
		},
		// Valid case: List second page with 3 employees per page
		{
			name:              "List second page with 3 employees per page",
			page:              2,
			size:              3,
			expectedCode:      http.StatusOK,
			expectedBodyCount: 3,
		},
		// Invalid page number: Page 0
		{
			name:         "Invalid page number: Page 0",
			page:         0,
			size:         2,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid page number",
		},
		// Invalid page size: Size 0
		{
			name:         "Invalid page size: Size 0",
			page:         1,
			size:         0,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid page size",
		},
		// No employees found
		{
			name:         "No employees found",
			page:         5,
			size:         2,
			expectedCode: http.StatusNotFound,
			expectedBody: "No employees found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/employees?page="+strconv.Itoa(tt.page)+"&size="+strconv.Itoa(tt.size), nil)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(ListEmployeesHandler)

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
			} else {
				// Check response body for successful list case
				var employees []models.Employee
				if err := json.Unmarshal(recorder.Body.Bytes(), &employees); err != nil {
					t.Fatalf("error unmarshalling response body: %v", err)
				}

				if len(employees) != tt.expectedBodyCount {
					t.Errorf("handler returned unexpected number of employees: got %v want %v",
						len(employees), tt.expectedBodyCount)
				}
			}
		})
	}
}

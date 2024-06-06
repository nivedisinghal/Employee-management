package handlers

import (
	"bytes"
	"ems/models"
	"ems/store"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestUpdateEmployeeHandler(t *testing.T) {
	// Initialize some employees for testing
	store.CreateEmployee("John Doe", "Developer", 60000.0)
	store.CreateEmployee("Alice Smith", "Manager", 80000.0)

	tests := []struct {
		name         string
		id           int
		payload      interface{}
		expectedCode int
		expectedBody string
	}{
		// Valid case: Update existing employee
		{
			name:         "Update existing employee",
			id:           1,
			payload:      models.Employee{Name: "Updated John Doe", Position: "Senior Developer", Salary: 70000.0},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"name":"Updated John Doe","position":"Senior Developer","salary":70000}`,
		},
		// Employee does not exist
		{
			name:         "Update non-existent employee",
			id:           80,
			payload:      models.Employee{Name: "New Employee", Position: "Tester", Salary: 50000.0},
			expectedCode: http.StatusNotFound,
			expectedBody: "Employee not found",
		},
		// Invalid payload: Missing name
		{
			name:         "Invalid payload: Missing name",
			id:           1,
			payload:      models.Employee{Position: "Senior Developer", Salary: 70000.0},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid employee data",
		},
		// Invalid payload: Negative salary
		{
			name:         "Invalid payload: Negative salary",
			id:           1,
			payload:      models.Employee{Name: "Updated John Doe", Position: "Senior Developer", Salary: -50000.0},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid employee data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("PUT", "/employees/"+strconv.Itoa(tt.id), bytes.NewBuffer(payload))
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(UpdateEmployeeHandler)

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
				// Check response body for successful update case
				var updatedEmployee models.Employee
				if err := json.Unmarshal(recorder.Body.Bytes(), &updatedEmployee); err != nil {
					t.Fatalf("error unmarshalling response body: %v", err)
				}

				if updatedEmployee.Name != tt.payload.(models.Employee).Name ||
					updatedEmployee.Position != tt.payload.(models.Employee).Position ||
					updatedEmployee.Salary != tt.payload.(models.Employee).Salary {
					t.Errorf("handler returned unexpected body: got %v want %v",
						updatedEmployee, tt.payload)
				}
			}
		})
	}
}

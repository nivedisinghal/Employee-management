package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ems/models"
)

func TestCreateEmployeeHandler(t *testing.T) {
	tests := []struct {
		name         string
		position     string
		salary       float64
		expectedCode int
	}{
		// Valid employee data
		{
			name:         "John Doe",
			position:     "Developer",
			salary:       60000.0,
			expectedCode: http.StatusCreated,
		},
		// Invalid data: empty name
		{
			name:         "",
			position:     "Manager",
			salary:       80000.0,
			expectedCode: http.StatusBadRequest,
		},
		// Invalid data: empty position
		{
			name:         "Alice Smith",
			position:     "",
			salary:       90000.0,
			expectedCode: http.StatusBadRequest,
		},
		// Invalid data: negative salary
		{
			name:         "Bob Johnson",
			position:     "Designer",
			salary:       -50000.0,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employee := models.Employee{
				Name:     tt.name,
				Position: tt.position,
				Salary:   tt.salary,
			}

			payload, _ := json.Marshal(employee)
			req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(CreateEmployeeHandler)

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					recorder.Code, tt.expectedCode)
			}

			// Check response body for valid cases
			if tt.expectedCode == http.StatusCreated {
				var createdEmployee models.Employee
				err := json.Unmarshal(recorder.Body.Bytes(), &createdEmployee)
				if err != nil {
					t.Fatalf("error unmarshalling response body: %v", err)
				}

				if createdEmployee.Name != tt.name || createdEmployee.Position != tt.position || createdEmployee.Salary != tt.salary {
					t.Errorf("handler returned unexpected body: got %v, want %v", createdEmployee, employee)
				}
			}
		})
	}

	// Additional test case: invalid JSON payload
	t.Run("Invalid JSON payload", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer([]byte("{")))
		if err != nil {
			t.Fatal(err)
		}
		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateEmployeeHandler)

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for invalid JSON: got %v want %v",
				recorder.Code, http.StatusBadRequest)
		}
	})
}

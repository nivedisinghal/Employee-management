package handlers

import (
	"bytes"
	"ems/store"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestDeleteEmployeeHandler(t *testing.T) {
	// Initialize some employees for testing
	store.CreateEmployee("Eve Williams", "Tester", 65000.0)
	tests := []struct {
		name         string
		id           int
		expectedCode int
		expectedBody string
	}{
		// Valid case: Delete existing employee
		{
			name:         "Valid case: Delete existing employee",
			id:           2,
			expectedCode: http.StatusNoContent,
			expectedBody: "Employee Deleted Successfully",
		},
		// Invalid case: Delete non-existing employee
		{
			name:         "Invalid case: Delete non-existing employee",
			id:           90,
			expectedCode: http.StatusNotFound,
			expectedBody: "Employee not found",
		},
		// Invalid case: Invalid employee ID
		{
			name:         "Invalid case: Invalid employee ID",
			id:           0,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid employee ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/employees/"+strconv.Itoa(tt.id), nil)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(DeleteEmployeeHandler)

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					recorder.Code, tt.expectedCode)
			}

			if tt.expectedBody != "" {
				body := bytes.TrimSpace(recorder.Body.Bytes())
				if !bytes.Contains(body, []byte(tt.expectedBody)) {
					t.Errorf("handler returned unexpected body: got %v want %v",
						string(body), tt.expectedBody)
				}
			}
		})
	}
}

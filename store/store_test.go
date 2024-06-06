package store

import (
	"ems/models"
	"errors"
	"reflect"
	"sync"
	"testing"
)

func TestCreateEmployee(t *testing.T) {
	tests := []struct {
		name     string
		position string
		salary   float64
		want     models.Employee
	}{
		{
			name:     "John Doe",
			position: "Developer",
			salary:   60000.0,
			want: models.Employee{
				ID:       1,
				Name:     "John Doe",
				Position: "Developer",
				Salary:   60000.0,
			},
		},
		{
			name:     "Alice Smith",
			position: "Manager",
			salary:   80000.0,
			want: models.Employee{
				ID:       2,
				Name:     "Alice Smith",
				Position: "Manager",
				Salary:   80000.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateEmployee(tt.name, tt.position, tt.salary)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateEmployee() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEmployeeByID(t *testing.T) {
	// Initialize some employees for testing
	employees[1] = models.Employee{ID: 1, Name: "John Doe", Position: "Developer", Salary: 60000.0}
	employees[2] = models.Employee{ID: 2, Name: "Alice Smith", Position: "Manager", Salary: 80000.0}

	tests := []struct {
		name        string
		id          int
		expected    models.Employee
		expectedErr error
	}{
		// Valid case: Employee exists
		{
			name:        "Valid employee",
			id:          1,
			expected:    models.Employee{ID: 1, Name: "John Doe", Position: "Developer", Salary: 60000.0},
			expectedErr: nil,
		},
		// Employee does not exist
		{
			name:        "Non-existent employee",
			id:          3,
			expected:    models.Employee{},
			expectedErr: errors.New("employee not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employee, err := GetEmployeeByID(tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("GetEmployeeByID() error = %v, want %v", err, tt.expectedErr)
				}
			} else {
				if err != nil {
					t.Errorf("GetEmployeeByID() unexpected error: %v", err)
				}
				if employee.ID != tt.expected.ID ||
					employee.Name != tt.expected.Name ||
					employee.Position != tt.expected.Position ||
					employee.Salary != tt.expected.Salary {
					t.Errorf("GetEmployeeByID() = %v, want %v", employee, tt.expected)
				}
			}
		})
	}
}

func TestUpdateEmployee(t *testing.T) {
	// Initialize some employees for testing
	CreateEmployee("John Doe", "Developer", 60000.0)
	CreateEmployee("Alice Smith", "Manager", 80000.0)

	tests := []struct {
		name           string
		id             int
		nameToUpdate   string
		positionUpdate string
		salaryUpdate   float64
		expectedName   string
		expectedError  error
	}{
		// Valid case: Employee exists, update successful
		{
			name:           "Update valid employee",
			id:             1,
			nameToUpdate:   "Updated John Doe",
			positionUpdate: "Senior Developer",
			salaryUpdate:   70000.0,
			expectedName:   "Updated John Doe",
			expectedError:  nil,
		},
		// Employee does not exist
		{
			name:           "Update non-existent employee",
			id:             19,
			nameToUpdate:   "New Employee",
			positionUpdate: "Tester",
			salaryUpdate:   50000.0,
			expectedName:   "",
			expectedError:  errors.New("Employee not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedEmployee, err := UpdateEmployee(tt.id, tt.nameToUpdate, tt.positionUpdate, tt.salaryUpdate)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("UpdateEmployee() error = %v, want %v", err, tt.expectedError)
				}
			} else {
				if err != nil {
					t.Errorf("UpdateEmployee() unexpected error: %v", err)
				}
				if updatedEmployee.Name != tt.expectedName {
					t.Errorf("UpdateEmployee() = %v, want %v", updatedEmployee.Name, tt.expectedName)
				}
			}
		})
	}
}

func TestListEmployees(t *testing.T) {
	// Initialize some employees for testing
	CreateEmployee("John Doe", "Developer", 60000.0)
	CreateEmployee("Alice Smith", "Manager", 80000.0)
	CreateEmployee("Bob Johnson", "Designer", 70000.0)
	CreateEmployee("Eve Williams", "Tester", 65000.0)

	tests := []struct {
		name          string
		page          int
		perPage       int
		expectedCount int
	}{
		// Valid case: List first page with 2 employees per page
		{
			name:          "List first page with 2 employees per page",
			page:          1,
			perPage:       2,
			expectedCount: 2,
		},
		// Valid case: List second page with 3 employees per page
		{
			name:          "List second page with 3 employees per page",
			page:          2,
			perPage:       3,
			expectedCount: 3,
		},
		// Valid case: List third page with 10 employees per page (only 2 employees available)
		{
			name:          "List third page with 10 employees per page (only 2 employees available)",
			page:          3,
			perPage:       10,
			expectedCount: 0,
		},
		// Invalid case: Page 0, perPage 2
		{
			name:          "Invalid case: Page 0, perPage 2",
			page:          0,
			perPage:       2,
			expectedCount: 0,
		},
		// Invalid case: Page 1, perPage 0
		{
			name:          "Invalid case: Page 1, perPage 0",
			page:          1,
			perPage:       0,
			expectedCount: 0,
		},
	}

	// Use a mutex to synchronize access to employees map
	var mu sync.Mutex

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Perform the test within a mutex lock to ensure no race conditions
			var paginatedEmployees []models.Employee
			mu.Lock()
			paginatedEmployees = ListEmployees(tt.page, tt.perPage)
			mu.Unlock()

			if len(paginatedEmployees) != tt.expectedCount {
				t.Errorf("ListEmployees returned unexpected number of employees: got %v want %v",
					len(paginatedEmployees), tt.expectedCount)
			}
		})
	}
}

func TestDeleteEmployee(t *testing.T) {

	tests := []struct {
		name         string
		id           int
		expectedErr  error
		expectedSize int
	}{
		// Valid case: Delete existing employee
		{
			name:         "Valid case: Delete existing employee",
			id:           1,
			expectedErr:  nil,
			expectedSize: 7,
		},
		// Invalid case: Delete non-existing employee
		{
			name:         "Invalid case: Delete non-existing employee",
			id:           20,
			expectedErr:  errors.New("Employee not found"),
			expectedSize: 7,
		},
	}

	// Use a mutex to synchronize access to employees map
	var mu sync.Mutex

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Perform the test within a mutex lock to ensure no race conditions
			mu.Lock()
			err := DeleteEmployee(tt.id)
			size := len(employees)
			mu.Unlock()

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("DeleteEmployee returned unexpected error: got %v, want %v", err, tt.expectedErr)
			}

			if size != tt.expectedSize {
				t.Errorf("DeleteEmployee did not delete the employee correctly: got %d, want %d", size, tt.expectedSize)
			}
		})
	}
}

package store

import (
	"ems/models"
	"errors"
	"sync"
)

var (
	employees = make(map[int]models.Employee)
	nextID    = 1
	mu        sync.Mutex
)

func CreateEmployee(name, position string, salary float64) models.Employee {
	mu.Lock()
	defer mu.Unlock()

	employee := models.Employee{
		ID:       nextID,
		Name:     name,
		Position: position,
		Salary:   salary,
	}
	employees[nextID] = employee
	nextID++

	return employee
}

func GetEmployeeByID(id int) (models.Employee, error) {
	mu.Lock()
	defer mu.Unlock()

	employee, exists := employees[id]
	if !exists {
		return models.Employee{}, errors.New("employee not found")
	}
	return employee, nil
}

func UpdateEmployee(id int, name, position string, salary float64) (models.Employee, error) {
	mu.Lock()
	defer mu.Unlock()

	employee, exists := employees[id]
	if !exists {
		return models.Employee{}, errors.New("Employee not found")
	}

	employee.Name = name
	employee.Position = position
	employee.Salary = salary
	employees[id] = employee

	return employee, nil
}

func DeleteEmployee(id int) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := employees[id]; !exists {
		return errors.New("Employee not found")
	}
	delete(employees, id)
	return nil
}

func ListEmployees(page, perPage int) []models.Employee {
	mu.Lock()
	defer mu.Unlock()

	var paginatedEmployees []models.Employee
	start := (page - 1) * perPage
	end := start + perPage

	count := 0
	for _, employee := range employees {
		if count >= start && count < end {
			paginatedEmployees = append(paginatedEmployees, employee)
		}
		count++
	}

	return paginatedEmployees
}

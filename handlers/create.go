package handlers

import (
	"ems/models"
	"ems/store"
	"encoding/json"
	"net/http"
)

func CreateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if employee.Name == "" || employee.Position == "" || employee.Salary <= 0 {
		http.Error(w, "Invalid employee data", http.StatusBadRequest)
		return
	}

	createdEmployee := store.CreateEmployee(employee.Name, employee.Position, employee.Salary)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEmployee)
}

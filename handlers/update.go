package handlers

import (
	"ems/models"
	"ems/store"
	"encoding/json"
	"net/http"
	"strconv"
)

func UpdateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/employees/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if employee.Name == "" || employee.Position == "" || employee.Salary <= 0 {
		http.Error(w, "Invalid employee data", http.StatusBadRequest)
		return
	}

	updatedEmployee, err := store.UpdateEmployee(id, employee.Name, employee.Position, employee.Salary)
	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedEmployee)
}

package handlers

import (
	"ems/store"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/employees/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	employee, err := store.GetEmployeeByID(id)
	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

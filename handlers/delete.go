package handlers

import (
	"ems/store"
	"net/http"
	"strconv"
)

func DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/employees/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	err = store.DeleteEmployee(id)
	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Employee Deleted Successfully"))
}

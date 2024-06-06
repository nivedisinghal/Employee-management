package handlers

import (
	"ems/store"
	"encoding/json"
	"net/http"
	"strconv"
)

func ListEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	perPage, err := strconv.Atoi(query.Get("size"))
	if err != nil || perPage < 1 {
		http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}

	employees := store.ListEmployees(page, perPage)
	if len(employees) == 0 {
		http.Error(w, "No employees found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

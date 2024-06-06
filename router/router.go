package router

import (
	"ems/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/employees", handlers.CreateEmployeeHandler).Methods("POST")

	return router
}

package main

import (
	"ems/router"
	"log"
	"net/http"
)

func main() {
	r := router.SetupRouter()
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}

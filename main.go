package main

import (
	"log"
	"net/http"
	"package_tracking_backend/handlers"

	"github.com/gorilla/mux"
)

func main() {
	ConnectDB()
	r := mux.NewRouter()
	r.HandleFunc("/api/users/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/api/users/login", handlers.LoginUser).Methods("POST")
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	controller "go-user-auth/controller"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/register", controller.RegisterUser).Methods("POST")
	router.HandleFunc("/getUser", controller.AllUser).Methods("GET")
	router.HandleFunc("/insertUser", controller.InsertAccount).Methods("POST")
	router.HandleFunc("/loginUser", controller.Login).Methods("POST")
	http.Handle("/", router)
	fmt.Println("Connected to port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

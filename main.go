package main

import (
	"database/sql"
	"fmt"
	"go-user-auth/handlers"
	"go-user-auth/middleware"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	controller "go-user-auth/controller"
)

var DB *sql.DB

// this functionn makes a variable DB to hold db connection, to access it throughout the project
func ConnectDB() {
	db, err := sql.Open("mysql", "root:@(localhost:3000)/userdb?parseTime=true")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("MYSQL is connected .....")
	DB = db

}

func main() {
	ConnectDB()
	router := mux.NewRouter()
	router.HandleFunc("/api/register", controller.RegisterUser).Methods("POST")
	router.HandleFunc("/getUser", controller.AllUser).Methods("GET")
	router.HandleFunc("/insertUser", controller.InsertAccount).Methods("POST")
	router.HandleFunc("/loginUser", controller.Login).Methods("POST")
	router.HandleFunc("/api/orders", handlers.CreateOrder).Methods("POST")
	//router.HandleFunc("/api/orders", middleware.AuthMiddleware(handlers.CreateOrder)).Methods("POST")
	router.HandleFunc("/api/orders/user/{user_id}", middleware.AuthMiddleware(handlers.GetOrders)).Methods("GET")
	router.HandleFunc("/api/orders/{id}", middleware.AuthMiddleware(handlers.GetOrderDetails)).Methods("GET")
	router.HandleFunc("/api/couriers/assigned_orders", middleware.AuthMiddleware(handlers.GetAssignedOrders)).Methods("GET")
	router.HandleFunc("/api/couriers/update_status", middleware.AuthMiddleware(handlers.UpdateOrderStatus)).Methods("PUT")
	router.HandleFunc("/api/admin/orders", middleware.AuthMiddleware(handlers.GetAllOrders)).Methods("GET")
	router.HandleFunc("/api/admin/assign_order", middleware.AuthMiddleware(handlers.AssignOrderToCourier)).Methods("POST")
	router.HandleFunc("/api/admin/delete_order/{order_id}", middleware.AuthMiddleware(handlers.DeleteOrder)).Methods("DELETE")
	http.Handle("/", router)

	log.Println("Received request to create order")
	fmt.Println("Connected to port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

// handlers are also controller functions, resposible for
// processing http requests for specific routes, ecxh handler fn
// is associated with a route and performs a specific action

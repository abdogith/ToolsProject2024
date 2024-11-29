package main

import (
	"log"
	"net/http"
	"package_tracking_backend/config"
	"package_tracking_backend/handlers"
	"package_tracking_backend/middleware"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	r := mux.NewRouter()

	r.HandleFunc("/testdb", func(w http.ResponseWriter, r *http.Request) {
		if config.DB != nil {
			w.Write([]byte("Database is initialized"))
		} else {
			http.Error(w, "Database is not initialized", http.StatusInternalServerError)
		}
	}).Methods("GET")
	r.HandleFunc("/api/users/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/api/users/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/api/orders", middleware.AuthMiddleware(handlers.CreateOrder)).Methods("POST")
	r.HandleFunc("/api/orders/user/{user_id}", middleware.AuthMiddleware(handlers.GetOrders)).Methods("GET")
	r.HandleFunc("/api/orders/{order_id}", middleware.AuthMiddleware(handlers.GetOrderDetails)).Methods("GET")
	r.HandleFunc("/api/couriers/assigned_orders/{courier_id}", middleware.AuthMiddleware(handlers.GetAssignedOrders)).Methods("GET")
	r.HandleFunc("/api/couriers/update_status", middleware.AuthMiddleware(handlers.UpdateOrderStatus)).Methods("PUT")
	r.HandleFunc("/api/admin/orders", middleware.AuthMiddleware(handlers.GetAllOrders)).Methods("GET")
	r.HandleFunc("/api/admin/assign_order", middleware.AuthMiddleware(handlers.AssignOrderToCourier)).Methods("POST")
	r.HandleFunc("/api/admin/delete_order/{order_id}", middleware.AuthMiddleware(handlers.DeleteOrder)).Methods("DELETE")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

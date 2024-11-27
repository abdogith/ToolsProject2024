package handlers

import (
	"database/sql"
	"encoding/json"
	"go-user-auth/middleware"
	"go-user-auth/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

// RegisterUser handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	// Decode the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert the new user into the database
	_, err = DB.Exec("INSERT INTO users (name, email, phone, password, role) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.Email, user.Phone, string(hashedPassword), "user")
	if err != nil {
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User registered successfully")
}

// LoginUser handles user login and returns a JWT token
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	// Decode the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Query the user by email from the database
	row := DB.QueryRow("SELECT id, password, role FROM users WHERE email = ?", user.Email)
	var dbUser model.User
	if err := row.Scan(&dbUser.Id, &dbUser.Password, &dbUser.Role); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token for the user
	token, err := middleware.GenerateToken(dbUser.Id, dbUser.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Respond with the JWT token
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// CreateOrder handles creating a new order
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	// Decode the request body into the order struct
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Println("Error decoding request body:", err) // Add logging for the error
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert the new order into the database
	res, err := DB.Exec("INSERT INTO orders (user_id, pickup_location, dropoff_location, package_details, delivery_time, status) VALUES (?, ?, ?, ?, ?, ?)",
		order.UserID, order.PickupLocation, order.DropoffLocation, order.PackageDetails, order.DeliveryTime, "pending")
	if err != nil {
		log.Println("Error creating order in the database:", err) // Add logging for the error
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Get the last inserted order ID
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error retrieving order ID:", err) // Add logging for the error
		http.Error(w, "Failed to retrieve order ID", http.StatusInternalServerError)
		return
	}

	// Set the order ID and respond with the created order
	order.ID = int(id)

	// Send a response with status 201 Created and the order details
	w.WriteHeader(http.StatusCreated) // Use status 201 for successful creation
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		log.Println("Error encoding response:", err) // Add logging for the error
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetOrders handles fetching orders for a user
func GetOrders(w http.ResponseWriter, r *http.Request) {
	// Fetch user ID from the URL
	userID := mux.Vars(r)["user_id"]

	// Query the database to fetch the orders for the given user
	rows, err := DB.Query("SELECT id, pickup_location, dropoff_location, package_details, delivery_time, status FROM orders WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Prepare the slice of orders to hold the results
	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.ID, &order.PickupLocation, &order.DropoffLocation, &order.PackageDetails, &order.DeliveryTime, &order.Status); err != nil {
			http.Error(w, "Failed to scan order", http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	// Check if there was any error during the iteration
	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	// Respond with the fetched orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetOrderDetails retrieves details of a specific order by ID
func GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var order model.Order
	// Query the order by ID
	err := DB.QueryRow("SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status FROM orders WHERE id = ?", id).
		Scan(&order.ID, &order.UserID, &order.PickupLocation, &order.DropoffLocation, &order.PackageDetails, &order.DeliveryTime, &order.Status)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	// Respond with the order details
	json.NewEncoder(w).Encode(order)
}

// GetAssignedOrders retrieves orders assigned to a specific courier
func GetAssignedOrders(w http.ResponseWriter, r *http.Request) {
	courierID := r.Context().Value("userID").(int) // Assuming userID is set by AuthMiddleware
	// Query the orders assigned to the courier
	rows, err := DB.Query("SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status FROM orders WHERE courier_id = ?", courierID)
	if err != nil {
		http.Error(w, "Failed to retrieve assigned orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect the orders
	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.PickupLocation, &order.DropoffLocation, &order.PackageDetails, &order.DeliveryTime, &order.Status); err != nil {
			http.Error(w, "Failed to retrieve assigned orders", http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	// Respond with the orders
	json.NewEncoder(w).Encode(orders)
}

// UpdateOrderStatus allows couriers to update the status of an order
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	var request struct {
		OrderID int    `json:"order_id"`
		Status  string `json:"status"`
	}
	// Decode the request payload
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the order status in the database
	_, err := DB.Exec("UPDATE orders SET status = ? WHERE id = ?", request.Status, request.OrderID)
	if err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	json.NewEncoder(w).Encode("Order status updated")
}

// GetAllOrders allows admin to retrieve all orders
func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	// Query all orders
	rows, err := DB.Query("SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status FROM orders")
	if err != nil {
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect the orders
	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.PickupLocation, &order.DropoffLocation, &order.PackageDetails, &order.DeliveryTime, &order.Status); err != nil {
			http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	// Respond with the orders
	json.NewEncoder(w).Encode(orders)
}

// AssignOrderToCourier allows admin to assign an order to a courier
func AssignOrderToCourier(w http.ResponseWriter, r *http.Request) {
	var request struct {
		OrderID   int `json:"order_id"`
		CourierID int `json:"courier_id"`
	}

	// Decode the request body into the request structure
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the order in the database with the new courier_id
	_, err := DB.Exec("UPDATE orders SET courier_id = ? WHERE id = ?", request.CourierID, request.OrderID)
	if err != nil {
		http.Error(w, "Failed to assign order to courier", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	json.NewEncoder(w).Encode("Order assigned to courier successfully")
}

// DeleteOrder allows admin to delete an order by ID
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["order_id"]

	// Delete the order from the database using the provided order ID
	_, err := DB.Exec("DELETE FROM orders WHERE id = ?", orderID)
	if err != nil {
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	json.NewEncoder(w).Encode("Order deleted successfully")
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"package_tracking_backend/config"
	"package_tracking_backend/middleware"
	"package_tracking_backend/models"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = config.DB.Exec("INSERT INTO users (name, email, phone, password, role) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.Email, user.Phone, string(hashedPassword), "user")
	if err != nil {
		log.Println("Error executing DB query:", err)
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User registered successfully")
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Parsed login request: email=%s", user.Email)

	// Query the database using user email
	row := config.DB.QueryRow("SELECT user_id, password, role FROM users WHERE email = ?", user.Email)
	var dbUser models.User
	err = row.Scan(&dbUser.UserID, &dbUser.Password, &dbUser.Role)
	if err == sql.ErrNoRows {
		// If no user is found, return error
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		// Handle other errors (database issues, etc.)
		log.Printf("Error scanning row: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Compare the provided password with the hashed password from the database
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		// If passwords don't match, return error
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token with user ID and role
	token, err := middleware.GenerateToken(dbUser.UserID, dbUser.Role)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Respond with both the token and a success message
	response := map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateOrder handles creating a new order
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	res, err := config.DB.Exec("INSERT INTO orders (user_id, pickup_location, dropoff_location, package_details, delivery_time, status) VALUES (?, ?, ?, ?, ?, ?)",
		order.UserID, order.PickupLocation, order.DropoffLocation, order.PackageDetails, order.DeliveryTime, "pending")
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	order.ID = int(id)
	json.NewEncoder(w).Encode(order)
}

// GetOrders retrieves all orders for a specific user
func GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	log.Printf("Fetching orders for user_id: %s", userID)

	rows, err := config.DB.Query("SELECT order_id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status FROM orders WHERE user_id = ?", userID)
	if err != nil {
		log.Printf("Database query failed: %v", err)
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.PickupLocation, &order.DropoffLocation, &order.PackageDetails, &order.DeliveryTime, &order.Status); err != nil {
			log.Printf("Row scanning failed: %v", err)
			http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		log.Printf("No orders found for user_id: %s", userID)
		http.Error(w, "No orders found", http.StatusNotFound)
		return
	}

	log.Printf("Orders retrieved successfully for user_id: %s", userID)
	json.NewEncoder(w).Encode(orders)
}

func GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["order_id"]
	var order models.Order
	err := config.DB.QueryRow("SELECT order_id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status FROM orders WHERE order_id = ?", id).
		Scan(&order.ID, &order.UserID, &order.PickupLocation, &order.DropoffLocation, &order.PackageDetails, &order.DeliveryTime, &order.Status)
	if err != nil {
		// Handle case where the order is not found or other errors occur
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func GetAssignedOrders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courierID, ok := vars["courier_id"]
	if !ok {
		log.Println("Courier ID not provided in URL")
		http.Error(w, "Courier ID is required", http.StatusBadRequest)
		return
	}

	// Convert courierID to integer
	courierIDInt, err := strconv.Atoi(courierID)
	if err != nil {
		log.Printf("Invalid courier ID: %v", err)
		http.Error(w, "Invalid courier ID", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching assigned orders for courier ID: %d", courierIDInt)

	// Query database for assigned orders for the specific courier
	rows, err := config.DB.Query(`
        SELECT o.order_id, o.user_id, o.pickup_location, o.dropoff_location, 
               o.package_details, o.delivery_time, o.status 
        FROM orders o
        JOIN assigned_orders ao ON o.order_id = ao.order_id
        WHERE ao.courier_id = ?`, courierIDInt)
	if err != nil {
		log.Printf("Database query failed: %v", err)
		http.Error(w, "Failed to retrieve assigned orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Initialize a slice to hold the orders
	var orders []models.Order

	// Iterate through the result set
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.PickupLocation, &order.DropoffLocation, &order.PackageDetails, &order.DeliveryTime, &order.Status); err != nil {
			log.Printf("Row scanning failed: %v", err)
			http.Error(w, "Failed to retrieve assigned orders", http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	// Check if no orders were found
	if len(orders) == 0 {
		log.Printf("No assigned orders found for courier ID: %d", courierIDInt)
		json.NewEncoder(w).Encode([]models.Order{}) // Return an empty array
		return
	}

	// Log success and respond with the list of orders
	log.Printf("Successfully retrieved %d assigned orders for courier ID: %d", len(orders), courierIDInt)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// UpdateOrderStatus allows couriers to update the status of an order
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	var request struct {
		OrderID int    `json:"order_id"`
		Status  string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := config.DB.Exec("UPDATE orders SET status = ? WHERE order_id = ?", request.Status, request.OrderID)
	if err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Order status updated")
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {

	query := "SELECT order_id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status FROM orders"

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Printf("Database query failed: %v", err)
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	count := 0

	for rows.Next() {
		log.Println("Fetching next row")
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.PickupLocation, &order.DropoffLocation, &order.PackageDetails, &order.DeliveryTime, &order.Status); err != nil {
			log.Printf("Row scanning failed: %v", err)
			http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
			return
		}
		log.Printf("Fetched order: %+v", order)
		orders = append(orders, order)
		count++
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		log.Println("No orders found")
		json.NewEncoder(w).Encode([]models.Order{}) // Return an empty array
		return
	}

	log.Printf("Successfully retrieved %d orders", count)
	json.NewEncoder(w).Encode(orders)
}

func AssignOrderToCourier(w http.ResponseWriter, r *http.Request) {
	var request struct {
		OrderID   int `json:"order_id"`
		CourierID int `json:"courier_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the order exists
	var orderExists int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE order_id = ?", request.OrderID).Scan(&orderExists)
	if err != nil || orderExists == 0 {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Check if the courier exists
	var courierExists int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM couriers WHERE courier_id = ?", request.CourierID).Scan(&courierExists)
	if err != nil || courierExists == 0 {
		http.Error(w, "Courier not found", http.StatusNotFound)
		return
	}

	// Remove existing assignment for the order (if any)
	_, err = config.DB.Exec("DELETE FROM assigned_orders WHERE order_id = ?", request.OrderID)
	if err != nil {
		http.Error(w, "Failed to clear existing assignment", http.StatusInternalServerError)
		return
	}

	// Insert the new assignment into the `assigned_orders` table
	_, err = config.DB.Exec("INSERT INTO assigned_orders (order_id, courier_id, status) VALUES (?, ?, ?)",
		request.OrderID, request.CourierID, "assigned")
	if err != nil {
		http.Error(w, "Failed to assign order to courier", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Order assigned to courier")
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["order_id"]

	// Check if the order exists
	var orderExists bool
	err := config.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE order_id = ?", orderID).Scan(&orderExists)
	if err != nil {
		http.Error(w, "Failed to check if order exists", http.StatusInternalServerError)
		return
	}
	if !orderExists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Delete related entries in the assigned_orders table
	_, err = config.DB.Exec("DELETE FROM assigned_orders WHERE order_id = ?", orderID)
	if err != nil {
		http.Error(w, "Failed to delete related order assignments", http.StatusInternalServerError)
		return
	}

	//  delete the order
	_, err = config.DB.Exec("DELETE FROM orders WHERE order_id = ?", orderID)
	if err != nil {
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Order deleted successfully")
}

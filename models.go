package models

type User struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type Order struct {
	ID              int    `json:"order_id"`
	UserID          int    `json:"user_id"`
	PickupLocation  string `json:"pickup_location"`
	DropoffLocation string `json:"dropoff_location"`
	PackageDetails  string `json:"package_details"`
	DeliveryTime    string `json:"delivery_time"`
	Status          string `json:"status"`
}

type Courier struct {
	ID       int    `json:"courier_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"-"`
}

type AssignedOrder struct {
	AssignmentID int    `json:"assignment_id"`
	OrderID      int    `json:"order_id"`
	CourierID    int    `json:"courier_id"`
	Status       string `json:"status"`
}

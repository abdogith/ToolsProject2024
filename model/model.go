package model

type User struct {
	Id       int    `form:"id" json:"id"`
	Name     string `form:"name" json:"name"`
	Email    string `form:"email" json:"email"`
	Phone    string `form:"phone" json:"phone"`
	Password string `form:"password" json:"password"`
	Role     string `json:"role"`
}
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User
}

type Order struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	PickupLocation  string `json:"pickup_location"`
	DropoffLocation string `json:"dropoff_location"`
	PackageDetails  string `json:"package_details"`
	DeliveryTime    string `json:"delivery_time"`
	Status          string `json:"status"`
}

type Courier struct {
	ID       int    `json:"id"`
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

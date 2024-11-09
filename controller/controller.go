package controller

import (
	"database/sql"
	"encoding/json"
	"go-user-auth/config"
	"go-user-auth/middleware"
	"go-user-auth/model"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func AllUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var response model.Response
	var arrUser []model.User

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT user_id, name, email, phone, password FROM users")
	if err != nil {
		log.Println("Error fetching users:", err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Ensure rows are closed after use

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Password)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue // Continue with the next row instead of terminating the program
		}
		arrUser = append(arrUser, user)
	}

	response.Status = 200
	response.Message = "Success"
	response.Data = arrUser

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func InsertAccount(w http.ResponseWriter, r *http.Request) {
	var response model.Response
	var newUser model.User

	// decoding the request body into the user struct
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return
	}

	db := config.Connect()
	defer db.Close()

	_, err = db.Exec("INSERT INTO users(name, email, phone, password) VALUES(?, ?, ?, ?)", newUser.Name, newUser.Email, newUser.Phone, newUser.Password)
	if err != nil {
		http.Error(w, "Error saving user to the database", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response.Status = 200
	response.Message = "Account inserted successfully"

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var response model.Response
	var newUser model.User

	// decoding the request body into the user struct
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return
	}

	// hashing the password before saving to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	db := config.Connect()
	defer db.Close()

	_, err = db.Exec("INSERT INTO users(name, email, phone, password) VALUES(?, ?, ?, ?)", newUser.Name, newUser.Email, newUser.Phone, newUser.Password)
	if err != nil {
		http.Error(w, "Error saving user to the database", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response.Status = 200
	response.Message = "User registered successfully"

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}
func Login(w http.ResponseWriter, r *http.Request) {
	var reqUser model.User
	err := json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := config.Connect()
	defer db.Close()

	var storedUser model.User
	log.Println("Executing query to fetch user by email:", reqUser.Email)
	// Update the query to use the correct column names: user_id and role
	err = db.QueryRow("SELECT user_id, name, email, phone, password, role FROM users WHERE email = ?", reqUser.Email).
		Scan(&storedUser.Id, &storedUser.Name, &storedUser.Email, &storedUser.Phone, &storedUser.Password, &storedUser.Role)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error querying user:", err) // Log the exact error
		http.Error(w, "Error querying user", http.StatusInternalServerError)
		return
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(reqUser.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := middleware.GenerateToken(storedUser.Id, storedUser.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Send the token in the response
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

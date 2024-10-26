package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"package_tracking_backend/middleware"
	"package_tracking_backend/models"

	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	_, err := DB.Exec("INSERT INTO users (name, email, phone, password, role) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.Email, user.Phone, string(hashedPassword), "user")
	if err != nil {
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User registered successfully")
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	row := DB.QueryRow("SELECT id, password, role FROM users WHERE email = ?", user.Email)
	var dbUser models.User
	row.Scan(&dbUser.ID, &dbUser.Password, &dbUser.Role)

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := middleware.GenerateToken(dbUser.ID, dbUser.Role)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

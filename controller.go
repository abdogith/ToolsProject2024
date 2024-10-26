package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-user-auth/config"
	"go-user-auth/model"

	"golang.org/x/crypto/bcrypt"
)

func AllUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var response model.Response
	var arrUser []model.User

	db := config.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, email, phone, password FROM users")
	if err != nil {
		log.Print(err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Password)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			arrUser = append(arrUser, user)
		}
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

	db := config.Connect()
	defer db.Close()

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	// Get values from the form
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")

	// Check if any of the fields are empty
	if name == "" || email == "" || phone == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Insert into the database
	_, err = db.Exec("INSERT INTO users(name, email, phone, password) VALUES(?, ?, ?, ?)", name, email, phone, password)
	if err != nil {
		http.Error(w, "Error saving user to the database", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response.Status = 200
	response.Message = "Account inserted successfully"
	fmt.Print("Insert account to database")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var response model.Response

	db := config.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users(name, email, phone, password) VALUES(?, ?, ?, ?)", name, email, phone, hashedPassword)
	if err != nil {
		http.Error(w, "Error saving user to the database", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response.Status = 200
	response.Message = "User registered successfully"
	fmt.Print("User registered successfully")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var reqUser model.User
	json.NewDecoder(r.Body).Decode(&reqUser)

	db := config.Connect()
	defer db.Close()

	var storedUser model.User
	err := db.QueryRow("SELECT id, name, email, phone, password FROM users WHERE email = ?", reqUser.Email).
		Scan(&storedUser.Id, &storedUser.Name, &storedUser.Email, &storedUser.Phone, &storedUser.Password)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error querying user", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(reqUser.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode("Login successful")
}

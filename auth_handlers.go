//4 august 2026 18:59

// func registerUser(w http.ResponseWriter, r *http.Request, db *sql.DB)
// func loginUser(w http.ResponseWriter, r *http.Request, db *sql.DB)

package main

import (
	// "encoding/json"
	// "strconv"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("very-secret-key")

func registerUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var input RegisterRequest

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &input)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	if len(input.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Problew with a bcrypt", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO users (name, age, email, password_hash) VALUES (?, ?, ?, ?)", input.Name, input.Age, input.Email, string(password_hash))

	if err != nil {
		http.Error(w, "Problew with a db.Exec", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "User registered")
}

func loginUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var input LoginRequest

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Problew with a ReadAll", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &input)

	if err != nil {
		http.Error(w, "Problew with a Unmarshaling", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	_, err = mail.ParseAddress(input.Email)

	if err != nil {
		http.Error(w, "Problew with an email format", http.StatusBadRequest)
		return
	}

	var user User

	err = db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", input.Email).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))

	if err != nil {
		http.Error(w, "Wrong email or password", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)

	type LoginResponse struct {
		Token string `json:"token"`
	}

	if err != nil {
		http.Error(w, "Problew with a JWT creating", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token: tokenString,
	}

	data, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Could not create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}

}

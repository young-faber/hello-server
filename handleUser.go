package main

import (
	// "encoding/json"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	_ "modernc.org/sqlite"
	"net/http"
)

// type User struct {
// 	ID   int
// 	Name string
// 	Age  int
// }

// var users []User
// var ID int

func handleUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &user)

	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := db.Exec(
		"INSERT INTO users (name, age) VALUES (?, ?)",
		user.Name,
		user.Age,
	)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	newID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not get user ID", http.StatusInternalServerError)
		return
	}

	user.ID = int(newID)

	fmt.Fprintf(w, "Added: %v.\n", user.Name)
}

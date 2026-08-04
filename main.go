// 26 july 2026 22:35

package main

import (
	// "encoding/json"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"net/http"
)

func main() {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL
	);
	`

	db, err := sql.Open("sqlite", "users.db")

	_, err = db.Exec(createTableQuery)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Hello, Web!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { //Function which prints message to client.
		fmt.Fprintf(w, "Hello, Web!")
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) { //Function which print message to client.
		name := r.URL.Query().Get("name")
		fmt.Printf("This is the about page. Method: %v. Path: %v\n Username: %v\n", r.Method, r.URL.Path, name)
		fmt.Fprintf(w, "This is the about page. Method: %v. Path: %v!\nYo! Wassup, %v ;)\n", r.Method, r.URL.Path, name)
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		handleUser(w, r, db)
	})

	err = http.ListenAndServe(":8000", nil)

	if err != nil {
		fmt.Println("Wrong connection to server.")
		// http.Error(w, "Error with connetction", http.StatusBadGateway)
	}

}

// 4 july 2026 finish.

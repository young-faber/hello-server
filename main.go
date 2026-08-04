// 26 july 2026 22:35

package main

import (
	// "encoding/json"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"net/http"
)

type User struct {
	ID   int
	Name string
	Age  int
}

var users []User
var ID int

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

	println(db)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

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
		switch r.Method {

		case http.MethodPost:
			CreateUser(w, r, db)
			return

		case http.MethodGet:
			ReadUser(w, r, db)
			return

		case http.MethodPut:
			UpdateUser(w, r, db)
			return

		case http.MethodDelete:
			DeleteUser(w, r, db)
			return

		default:
			http.Error(w, "Not allowed method", http.StatusMethodNotAllowed)
		}
	})

	err = http.ListenAndServe(":8000", nil)

	if err != nil {
		fmt.Println("Wrong connection to server.")
		// http.Error(w, "Error with connetction", http.StatusBadGateway)
	}

}

// 4 july 2026 finish.

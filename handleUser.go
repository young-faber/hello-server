package main

import (
	// "encoding/json"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"
)

// type User struct {
// 	ID   int
// 	Name string
// 	Age  int
// }

// var users []User
// var ID int

func CreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

func ReadUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var user User
	w.Header().Set("Content-Type", "application/json")
	idQuery := r.URL.Query().Get("id")
	if idQuery != "" {
		parsedID, err := strconv.Atoi(idQuery)

		if err != nil {
			fmt.Println("You have a problem with a strconv")
			http.Error(w, "Wrong input, ID has to be an integer", http.StatusBadRequest)
			return
		}

		err = db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", parsedID).Scan(&user.ID, &user.Name, &user.Age)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Problem with a DB.", http.StatusBadGateway)
			return
		}

		data, err := json.Marshal(user)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Problem with a Marshaling", http.StatusBadGateway)
			return
		}

		_, err = w.Write(data)
		return
		// if err != nil {
		// 	fmt.Println(err)
		// 	http.Error(w, "Problem with a writing", http.StatusBadGateway)
		// 	return
		// }

	} else {

		var users []User

		rows, err := db.Query("SELECT id, name, age FROM users ")

		if err != nil {
			fmt.Print(err)
			http.Error(w, "Problem with a DB", http.StatusBadGateway)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Age)

			if err != nil {
				fmt.Print(err)
				http.Error(w, "Wrong database request", http.StatusBadRequest)
				return
			}

			users = append(users, user)
			fmt.Println(user)
		}
		err = rows.Err()

		if err != nil {
			fmt.Print(err)
			http.Error(w, "Problem with a DB", http.StatusBadGateway)
			return
		}

		data, err := json.Marshal(users)
		var xxx int
		xxx, err = w.Write(data)
		fmt.Println(xxx, xxx, xxx)
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	idQuery := r.URL.Query().Get("id")

	user := User{}

	if idQuery != "" {
		parsedID, err := strconv.Atoi(idQuery)

		if err != nil {
			fmt.Println("You have a problem with a strconv")
			http.Error(w, "Wrong input, ID has to be an integer", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Println("You have a problem with a ReadAll")
			http.Error(w, "Wrong input, ID has to be an integer", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &user)

		if err != nil {
			fmt.Println("You have a problem with a Unmarshal")
			http.Error(w, "Wrong input, ID has to be an integer", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", user.Name, user.Age, parsedID)

		if err != nil {
			fmt.Println("You have a problem with a db.Exec")
			http.Error(w, "Wrong input, ID has to be an integer", http.StatusBadRequest)
			return
		}

		affected, err := result.RowsAffected()

		if affected == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if err != nil {
			fmt.Println("You have a problem with a RowsAffected")
			http.Error(w, "Wrong input, ID has to be an integer", http.StatusBadRequest)
			return
		}

		fmt.Printf("User with ID = %v, has been changed ", parsedID)
		fmt.Fprintf(w, "User with ID = %v has been changed", parsedID)

	} else {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	idQuery := r.URL.Query().Get("id")

	if idQuery != "" {
		parsedID, err := strconv.Atoi(idQuery)

		if err != nil {
			fmt.Println("You have a problem with a strconv")
			http.Error(w, "Wrong input, ID has to be an integer", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("DELETE FROM users WHERE id = ?", parsedID)

		if err != nil {
			fmt.Println("You have a problem with a DB")
			http.Error(w, "Problem with a DB", http.StatusBadRequest)
			return
		}

		affected, err := result.RowsAffected()

		if err != nil {
			fmt.Println("You have a problem with a Affected")
			http.Error(w, "Problem with a Affected", http.StatusBadRequest)
			return
		}

		if affected == 0 {
			http.Error(w, "There is not user with a such ID.", http.StatusNotFound)
			return
		}

		fmt.Printf("User with ID = %v has been deleted ", parsedID)
		fmt.Fprintf(w, "User with ID = %v has been deleted", parsedID)
		return

	} else {
		fmt.Fprint(w, "Write user's id.")
		fmt.Printf("User hasn't write ID.")
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
}

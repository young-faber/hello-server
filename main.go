// 26 july 2026 22:35

package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	// "path"
	"io"
)

type User struct {
	Id   int
	Name string
	Age  int
}

var users []User
var id int

func main() {
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
		var user User

		if r.Method == http.MethodPost {
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

			user.Id = id
			id++

			users = append(users, user)

			fmt.Fprintf(w, "Added: %v.\n", user.Name)
			return
		}
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			id_query := r.URL.Query().Get("id")
			if id_query != "" {
				new_id_query, err := strconv.Atoi(id_query)

				if err != nil {
					fmt.Println("You have a problem with a strconv")
					http.Error(w, "Wrong input, id has to be an integer", http.StatusBadRequest)
				}

				for _, user := range users {
					if user.Id == new_id_query {
						data, err := json.Marshal(user)

						if err != nil {
							fmt.Println("U have a problem with a Marshalization of user.")
						}

						w.Write(data)
						return
					}
				}
				http.Error(w, "There is not user with a such ID.", http.StatusNotFound)

			}

			data, err := json.Marshal(users)
			if err != nil {
				fmt.Printf("U have a problem, guys. %v", err)
				return
			}
			_, err = w.Write(data)
			return
		}
		http.Error(w, "Not allowed method", http.StatusMethodNotAllowed)

	})

	http.ListenAndServe(":8000", nil)

}

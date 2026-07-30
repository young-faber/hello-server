// 26 july 2026 22:35

package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"net/http"

	// "path"
	"io"
)

type User struct {
	Name string
	Age  int
}

var users []User

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

			users = append(users, user)

			fmt.Fprintf(w, "Added: %v.\n", user.Name)
		}
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			data, err := json.Marshal(users)
			if err != nil {
				fmt.Printf("U have some problem, guys. %v", err)
				return
			}
			_, err = w.Write(data)
		}

	})

	http.ListenAndServe(":8000", nil)

}

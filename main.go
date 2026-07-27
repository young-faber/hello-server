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

func main() {
	fmt.Println("Hello, Web!")

	type User struct {
		Name string
		Age  int
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Web!")
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		fmt.Printf("This is the about page. Method: %v. Path: %v\n Username: %v\n", r.Method, r.URL.Path, name)
		fmt.Fprintf(w, "This is the about page. Method: %v. Path: %v!\nYo! Wassup, %v ;)\n", r.Method, r.URL.Path, name)
	})
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		var user User
		// user := User{
		// 	Name: "Ars",
		// 	Age:  16,
		// }
		if r.Method == "GET" {
			fmt.Fprint(w, r.Method)
		} else {
			fmt.Fprint(w, r.Method)
		}

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

		fmt.Println(user, user.Age, user.Name)
	})

	http.ListenAndServe(":8000", nil)

}

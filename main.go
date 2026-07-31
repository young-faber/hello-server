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

	"golang.org/x/tools/go/analysis/passes/nilfunc"
)

type User struct {
	Id   int
	Name string
	Age  int
}

var users []User
var ID int

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

		switch r.Method {

		case http.MethodPost:
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

			ID++
			user.Id = ID

			users = append(users, user)

			fmt.Fprintf(w, "Added: %v.\n", user.Name)
			return

		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			idQuery := r.URL.Query().Get("id")
			if idQuery != "" {
				parsedID, err := strconv.Atoi(idQuery)

				if err != nil {
					fmt.Println("You have a problem with a strconv")
					http.Error(w, "Wrong input, ID has to be an integer", http.StatusBadRequest)
					return
				}

				for _, user := range users {
					if user.Id == parsedID {
						data, err := json.Marshal(user)

						if err != nil {
							fmt.Println("U have a problem with a Marshalization of user.")
							return
						}

						w.Write(data)
						return
					}
				}
				http.Error(w, "There is not user with a such ID.", http.StatusNotFound)

			} else {
				data, err := json.Marshal(users)
				if err != nil {
					fmt.Printf("U have a problem, guys. %v", err)
					return
				}
				_, err = w.Write(data)
				return
			}

		case http.MethodDelete:
			w.Header().Set("Content-Type", "application/json")
			idQuery := r.URL.Query().Get("id")

			if idQuery != "" {
				parsedID, err := strconv.Atoi(idQuery)

				if err != nil {
					fmt.Println("You have a problem with a strconv")
					http.Error(w, "Wrong input, ID has to be an integer", http.StatusBadRequest)
					return
				}

				for num, user := range users {
					if user.Id == parsedID {

						users = append(users[:num], users[num+1:]...)
						fmt.Printf("User with ID = %v, number %v, has been deleted ", user.Id, num)
						fmt.Fprintf(w, "User with ID = %v has been deleted", user.Id)
						return
					}
				}
				http.Error(w, "There is not user with a such ID.", http.StatusNotFound)
			} else {
				fmt.Fprint(w, "Write user's id.")
				fmt.Printf("User hasn't write ID.")
				return
			}

		case http.MethodPut:

			w.Header().Set("Content-Type", "application/json")
			idQuery := r.URL.Query().Get("id")

			if idQuery != "" {
				parsedID, err := strconv.Atoi(idQuery)

				if err != nil {
					fmt.Println("You have a problem with a strconv")
					http.Error(w, "Wrong input, ID has to be an integer", http.StatusBadRequest)
					return
				}

				for num, user := range users {
					if user.Id == parsedID {
						oldId := user.Id
						body, err := io.ReadAll(r.Body)

						if err != nil {
							fmt.Println(err)
							http.Error(w, "Wrong parameters.", http.StatusBadRequest)
							return
						}

						err = json.Unmarshal(body, &user)
						user.Id = oldId
						users[num] = user

						if err != nil {
							fmt.Println(err)
							http.Error(w, "Unmarshal wrong", http.StatusBadRequest)
							return
						}

						fmt.Printf("User with ID = %v, number %v, has been changed ", user.Id, num)
						fmt.Fprintf(w, "User with ID = %v has been changed", user.Id)
						return
					}
				}

				http.Error(w, "There is not user with a such ID.", http.StatusNotFound)
			} else {
				fmt.Fprint(w, "Write user's id.")
				fmt.Printf("User hasn't write ID")
				return
			}

		default:
			http.Error(w, "Not allowed method", http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		fmt.Println("Wrong connection to server.")
		// http.Error(w, "Error with connetction", http.StatusBadGateway)
	}

}

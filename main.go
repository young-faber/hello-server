// 26 july 2026 22:35

package main

import (
	"fmt"
	"net/http"
	// "path"
)

func main() {
	fmt.Println("Hello, Web!")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Web!")
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("This is the about page. Method: %v. Path: %v\n", r.Method, r.URL.Path)
	})
	http.ListenAndServe(":8080", nil)
}

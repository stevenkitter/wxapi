package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	fmt.Println("listening to 8021")
	http.ListenAndServe(":8021", nil)
}

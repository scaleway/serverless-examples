package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(string(body))
	fmt.Fprintln(w, "Hello from container")
}

func main() {
	http.HandleFunc("/", handle)

	log.Println("Starting!")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type HelloResponse struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/helloworld", helloHandler)
	http.HandleFunc("/versionz", versionz)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/helloworld" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello Stranger")
}

func versionz(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	response := HelloResponse{
		Message: fmt.Sprintf("Hello %s!", name),
	}
	jsonResponse(w, response, http.StatusOK)
}

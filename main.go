package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/helloworld", helloHandler) //simple router

	fmt.Printf("Starting server at port 8080, navigate to http://localhost:8080/helloworld\n") //set the message
	log.Println("Starting HTTP server on :8080")
	err := http.ListenAndServe(":8080", nil) //set the port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/helloworld" {
		fmt.Fprintf(w, "Hello Stranger")
	}
	if r.URL.Path != "/helloworld" {
		http.Error(w, "404, page not found.", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
}

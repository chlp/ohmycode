package main

import (
	"api/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/example", handlers.ExampleHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

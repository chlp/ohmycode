package main

import (
	"log"
	"net/http"
	"ohmycode_api/internal/api"
)

func main() {
	http.HandleFunc("/example", api.ExampleHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

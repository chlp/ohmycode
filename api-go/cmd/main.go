package main

import (
	"log"
	"net/http"
	"ohmycode_api/config"
	"ohmycode_api/internal/api"
	"ohmycode_api/internal/store"
)

func main() {
	http.HandleFunc("/example", api.ExampleHandler)
	apiConfig := config.LoadApiConf()
	db := store.NewDb(apiConfig.DB)
	v, err := db.Select("files", nil)
	println(len(v), v[0], err)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

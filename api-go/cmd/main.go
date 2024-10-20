package main

import (
	"log"
	"net/http"
	"ohmycode_api/config"
	"ohmycode_api/internal/api"
	"ohmycode_api/internal/store"
)

func main() {
	apiConfig := config.LoadApiConf()
	db := store.NewDb(apiConfig.DB)
	s := api.NewService(db)

	//v, err := db.Select("files", map[string]interface{}{"name": "abc"})
	//println(len(v), v[0], err)

	http.HandleFunc("/example", api.ExampleHandler)
	http.HandleFunc("/file/get_update", s.HandleFileGetUpdateRequest)
	//http.HandleFunc("/session", api.HandleSessionRequest)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

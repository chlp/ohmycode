package main

import (
	"ohmycode_api/config"
	"ohmycode_api/internal/api"
	"ohmycode_api/internal/store"
)

func main() {
	apiConfig := config.LoadApiConf()
	api.Run(store.NewStore(apiConfig.DB))
}

package config

import (
	"encoding/json"
	"log"
	"ohmycode_api/internal/store"
	"os"
)

const confPath = "api-conf.json"

type ApiConfig struct {
	DB store.DBConfig `json:"db"`
}

var conf ApiConfig

func LoadApiConf() ApiConfig {
	data, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatal("config: cannot read file")
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("config: cannot parse file")
	}
	return conf
}

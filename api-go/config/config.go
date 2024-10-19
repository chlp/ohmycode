package config

import (
	"encoding/json"
	"log"
	"ohmycode_api/internal/store"
	"os"
)

const confPath = "../api-config.json"

type Config struct {
	DB store.DBConfig `json:"store"`
}

var conf Config

func LoadApiConf() Config {
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

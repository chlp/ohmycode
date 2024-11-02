package config

import (
	"encoding/json"
	"log"
	"ohmycode_api/internal/store"
	"os"
)

const confPath = "api-conf.json"
const confExamplePath = "api-conf-example.json"

type ApiConfig struct {
	DB               store.DBConfig `json:"db"`
	HttpPort         int            `json:"http_port"`
	ServeClientFiles bool           `json:"serve_client_files"`
}

var conf ApiConfig

func LoadApiConf() ApiConfig {
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return loadConfFromFile(confExamplePath)
	}
	return loadConfFromFile(confPath)
}

func loadConfFromFile(filePath string) ApiConfig {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("config: cannot read file")
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("config: cannot parse file")
	}
	return conf
}

package conf

import (
	"encoding/json"
	"log"
	"os"
)

const confPath = "../api-conf.json"

type DBConfig struct {
	ServerName string `json:"servername"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	DBName     string `json:"dbname"`
}

type Config struct {
	DB DBConfig `json:"db"`
}

var conf Config

func LoadApiConf() Config {
	if conf.DB.ServerName != "" {
		return conf
	}

	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		log.Fatal("conf: please create conf file")
	}

	data, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatal("conf: cannot read file")
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("conf: cannot parse file")
	}

	return conf
}

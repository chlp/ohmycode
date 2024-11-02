package config

import (
	"encoding/json"
	"log"
	"ohmycode_runner/pkg/util"
	"os"
)

const confPath = "runner-conf.json"
const confExamplePath = "runner-conf-example.json"

type RunnerConf struct {
	RunnerId   string   `json:"id"`
	IsPublic   bool     `json:"is_public"`
	RunnerName string   `json:"name"`
	ApiUrl     string   `json:"api"`
	Languages  []string `json:"languages"`
}

var conf RunnerConf

func LoadRunnerConf() RunnerConf {
	var conf RunnerConf
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		conf = loadConfFromFile(confExamplePath)
		conf.RunnerId = util.GenUuid()
	} else {
		conf = loadConfFromFile(confPath)
	}
	if !util.IsUuid(conf.RunnerId) {
		log.Fatalf("config: runner id is wrong")
	}
	return conf
}

func loadConfFromFile(filePath string) RunnerConf {
	data, err := os.ReadFile(filePath)
	path, _ := os.Getwd()
	if err != nil {
		log.Fatalf("config: cannot read file: %s/%s", path, filePath)
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("config: cannot parse file: %s/%s", path, filePath)
	}
	return conf
}

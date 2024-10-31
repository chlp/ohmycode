package config

import (
	"encoding/json"
	"log"
	"ohmycode_runner/pkg/util"
	"os"
)

const confPath = "conf.json"

type RunnerConf struct {
	RunnerId   string   `json:"id"`
	IsPublic   bool     `json:"is_public"`
	RunnerName string   `json:"name"`
	ApiUrl     string   `json:"api"`
	Languages  []string `json:"languages"`
}

var conf RunnerConf

func LoadRunnerConf() RunnerConf {
	println(os.Getwd())
	data, err := os.ReadFile(confPath)
	path, _ := os.Getwd()
	if err != nil {
		log.Fatalf("config: cannot read file: %s/%s", path, confPath)
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("config: cannot parse file: %s/%s", path, confPath)
	}
	if !util.IsUuid(conf.RunnerId) {
		log.Fatalf("config: runner id is wrong: %s/%s", path, confPath)
	}
	return conf
}

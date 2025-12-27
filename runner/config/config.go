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

func LoadRunnerConf() RunnerConf {
	var conf RunnerConf
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		conf = loadConfFromFile(confExamplePath)
		conf.RunnerId = util.GenUuid()
		// Persist the generated runner id, otherwise it changes on each start.
		saveConfToFile(confPath, conf)
	} else {
		conf = loadConfFromFile(confPath)
	}
	if !util.IsUuid(conf.RunnerId) {
		log.Fatalf("config: runner id is wrong")
	}
	if conf.ApiUrl == "" {
		log.Fatalf("config: api url is empty")
	}
	return conf
}

func loadConfFromFile(filePath string) RunnerConf {
	data, err := os.ReadFile(filePath)
	path, _ := os.Getwd()
	if err != nil {
		log.Fatalf("config: cannot read file: %s/%s", path, filePath)
	}
	var conf RunnerConf
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("config: cannot parse file: %s/%s", path, filePath)
	}
	return conf
}

func saveConfToFile(filePath string, conf RunnerConf) {
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		log.Fatalf("config: cannot marshal config")
	}
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		log.Fatalf("config: cannot write file: %s", filePath)
	}
}

package config

import (
	"encoding/json"
	"log"
	"ohmycode_api/internal/store"
	"os"
	"strconv"
	"strings"
)

const confPath = "api-conf.json"
const confExamplePath = "api-conf-example.json"

type ApiConfig struct {
	DB               store.DBConfig `json:"db"`
	HttpPort         int            `json:"http_port"`
	ServeClientFiles bool           `json:"serve_client_files"`
	UseDynamicFiles  bool           `json:"use_dynamic_files"`
	// WsAllowedOrigins controls the Origin-check during WebSocket upgrade.
	// Empty (default) or ["*"] means allow all origins.
	// Otherwise it should contain allowed origins like "https://example.com" or hosts like "example.com:3000".
	WsAllowedOrigins []string `json:"ws_allowed_origins"`
	// ContentMaxLengthKb sets the maximum file content size in kilobytes. Default: 512.
	ContentMaxLengthKb int `json:"content_max_length_kb"`
}

func LoadApiConf() ApiConfig {
	var c ApiConfig
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		c = loadConfFromFile(confExamplePath)
	} else {
		c = loadConfFromFile(confPath)
	}
	applyEnvOverrides(&c)
	validateConf(c)
	return c
}

func applyEnvOverrides(c *ApiConfig) {
	if v := os.Getenv("OHMYCODE_MONGO_URI"); v != "" {
		c.DB.ConnectionString = v
	}
	if v := os.Getenv("OHMYCODE_MONGO_DBNAME"); v != "" {
		c.DB.DBName = v
	}
	if v := os.Getenv("OHMYCODE_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.HttpPort = port
		}
	}
	if v := os.Getenv("OHMYCODE_WS_ORIGINS"); v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		c.WsAllowedOrigins = parts
	}
}

func validateConf(c ApiConfig) {
	if c.DB.ConnectionString == "" {
		log.Fatal("config: OHMYCODE_MONGO_URI or db.connectionString is required")
	}
	if c.HttpPort == 0 {
		log.Fatal("config: OHMYCODE_PORT or http_port is required")
	}
}

func loadConfFromFile(filePath string) ApiConfig {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("config: cannot read file")
	}
	var c ApiConfig
	if err = json.Unmarshal(data, &c); err != nil {
		log.Fatal("config: cannot parse file")
	}
	return c
}

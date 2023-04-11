package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	DatabaseDriver string       `json:"database_driver"`
	DatabaseDSN    string       `json:"database_dsn"`
	Server         ServerConfig `json:"server"`
}

type ServerConfig struct {
	Port string `json:"port"`
}

var (
	// Instance of Config struct, accessible through the package
	Instance Config
)

// LoadConfig reads the configuration file and unmarshals it into Config struct.
func LoadConfig() {
	// Set default values
	Instance.Server.Port = "8080"

	// Load the configuration file
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Printf("Unable to open config file: %v, using defaults.", err)
		return
	}
	defer configFile.Close()

	// Read the configuration file
	bytes, err := io.ReadAll(configFile)
	if err != nil {
		log.Printf("Unable to read config file: %v, using defaults.", err)
		return
	}

	// Unmarshal the JSON configuration into Config struct
	err = json.Unmarshal(bytes, &Instance)
	if err != nil {
		log.Printf("Unable to parse config file: %v, using defaults.", err)
		return
	}
}

package types

import (
	"encoding/json"
	"os"
)

type Config struct {
	Addr     string `json:"Addr"`
	LogLevel string `json:"LogLevel"`
	LogFile  string `json:"LogFile"`
	MongoURI string `json:"MongoURI"`
	MongoDB  string `json:"MongoDB"`
}

func LoadConfig(configPath string) (*Config, error) {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

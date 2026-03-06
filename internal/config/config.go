package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL    string `json:"db_url"`
	Username string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func (c *Config) SetUser(userName string) {
	c.Username = userName
	write(*c)
}

func write(cfg Config) error {
	//write to config
	path, err := getConfigFilePath(configFileName)
	if err != nil {
		return fmt.Errorf("Path could not be resolved: %w", err)
	}
	jsonBytes, err := json.Marshal(cfg)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error: Failed to overwrite file: %w", err)
	}
	defer f.Close()
	f.Write(jsonBytes)
	return nil
}

func Read() (Config, error) {
	cfg := Config{}
	path, err := getConfigFilePath(configFileName)
	if err != nil {
		return Config{}, fmt.Errorf("Path could not be resolved: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, nil
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("Error unmarshaling data from config file: %w", err)
	}

	return cfg, nil
}

func getConfigFilePath(configPath string) (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := fmt.Sprintf("%s/%s", homePath, configPath)
	return fullPath, nil
}

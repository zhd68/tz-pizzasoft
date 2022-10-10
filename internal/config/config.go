package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	LogLavel        string `json:"log_lavel"`
	XAuthKey        string `json:"x_auth_key"`
	Listen          Listen `json:"listen"`
	StorageInMemory bool   `json:"storage_in_memory"`
	DB              DB     `json:"db"`
}

type Listen struct {
	BindIP string `json:"bind_ip"`
	Port   string `json:"port"`
}

type DB struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

func New() *Config {
	return &Config{
		LogLavel: "debug",
		XAuthKey: "qwerty123",
		Listen: Listen{
			BindIP: "localhost",
			Port:   "8080",
		},
		StorageInMemory: true,
	}
}

func (c *Config) ParseConfig(configPath string) error {
	filename, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("can't get config file: %w", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("can't read file: %w", err)
	}

	if err = json.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("can't unmarshal config: %w", err)
	}
	return nil
}

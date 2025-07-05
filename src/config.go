package main

import "os"
import (
	"encoding/json"
	"fmt"
)

type Config struct {
	rssUrl string
	contactsVcfFile string
	familyInfoFile string
	bedtimeHour int
	bedtimeMinute int
}

func (c *Config) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Config) FromJSON(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c *Config) SaveToFile(filename string) error {
	data, err := c.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	return os.WriteFile(filename, data, 0644)
}

func (c *Config) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	return c.FromJSON(data)
}

func getConfig() *Config {
	return &Config{
		Queue: os.Getenv("KIDSH_QUEUE"),
	}
}
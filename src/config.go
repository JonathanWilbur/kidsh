package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ConfigFile struct {
	RssUrl              string   `json:"rssUrl"`
	ContactsVcfFile     string   `json:"contactsVcfFile"`
	FamilyInfoFile      string   `json:"familyInfoFile"`
	BedtimeHour         int      `json:"bedtimeHour"`
	BedtimeMinute       int      `json:"bedtimeMinute"`
	BlacklistedCommands []string `json:"blacklistedCommands"`
	TodoFile            string   `json:"todoFile"`
	MyName              string   `json:"myName"`
	WeatherUrl          string   `json:"weatherUrl"`
}

type Config struct {
	RssUrl              string
	ContactsVcfFile     string
	FamilyInfoFile      string
	BedtimeHour         int
	BedtimeMinute       int
	BlacklistedCommands map[string]bool
	TodoFile            string
	MyName              string
	WeatherUrl          string
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

func configSearchPaths() []string {
	var paths []string
	if p := strings.TrimSpace(os.Getenv("KIDSH_CONFIG")); p != "" {
		paths = append(paths, p)
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		paths = append(paths, filepath.Join(xdg, "kidsh", "config.json"))
	} else if home, err := os.UserHomeDir(); err == nil && home != "" {
		paths = append(paths, filepath.Join(home, ".config", "kidsh", "config.json"))
	}
	if snap := os.Getenv("SNAP"); snap != "" {
		paths = append(paths, filepath.Join(snap, "etc", "kidsh.json"))
	}
	paths = append(paths, "/etc/kidsh.json", "/etc/kidsh/config.json")
	return paths
}

func getConfig() *Config {
	paths := configSearchPaths()
	var cf ConfigFile
	var found bool
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil {
			if err := json.Unmarshal(data, &cf); err == nil {
				found = true
				break
			}
		}
	}
	if !found {
		// If not found or failed to parse, return default config
		return &Config{
			RssUrl:              "",
			ContactsVcfFile:     "",
			FamilyInfoFile:      "",
			BedtimeHour:         21,
			BedtimeMinute:       0,
			BlacklistedCommands: make(map[string]bool),
			TodoFile:            "",
			MyName:              "",
			WeatherUrl:          "",
		}
	}
	blacklist := make(map[string]bool)
	for _, cmd := range cf.BlacklistedCommands {
		blacklist[cmd] = true
	}
	return &Config{
		RssUrl:              cf.RssUrl,
		ContactsVcfFile:     cf.ContactsVcfFile,
		FamilyInfoFile:      cf.FamilyInfoFile,
		BedtimeHour:         cf.BedtimeHour,
		BedtimeMinute:       cf.BedtimeMinute,
		BlacklistedCommands: blacklist,
		TodoFile:            cf.TodoFile,
		MyName:              cf.MyName,
		WeatherUrl:          cf.WeatherUrl,
	}
}

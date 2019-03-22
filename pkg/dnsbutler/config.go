package dnsbutler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Wait       int      `json:"waitInSec"`
	Provider   string   `json:"ipProvider"`
	ListenAddr string   `json:"listenAddr"`
	Targets    []string `json:"targets"`
}

func initConfig(configPath string) (*Config, error) {
	config := &Config{
		Wait:       5,
		Provider:   "https://api.ipify.org/",
		ListenAddr: ":5000",
		Targets:    make([]string, 0),
	}

	c, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(configPath, c, 0644)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func readConfig(configPath string) (*Config, error) {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func readOrInitConfig(configPath string, logger *log.Logger) (*Config, error) {
	var c *Config
	var err error

	if c, err = readConfig(configPath); err != nil && os.IsNotExist(err) {
		logger.Println("No config file found - generating default config file")
		c, err = initConfig(configPath)
	}

	if err != nil {
		return nil, err
	}

	return c, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/northbright/pathhelper"
	"github.com/shchnmz/ming"
	"github.com/shchnmz/zb"
)

// Config represents app settings.
type Config struct {
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

var (
	currentDir, configFile string
	config                 Config
)

func main() {
	var (
		err       error
		blacklist zb.Blacklist
	)

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	if err = loadConfig(configFile, &config); err != nil {
		return
	}

	db := zb.DB{ming.DB{config.RedisServer, config.RedisPassword}}
	if err = db.LoadBlacklistFromJSON("blacklist.json", &blacklist); err != nil {
		return
	}
}

// init initializes path variables.
func init() {
	currentDir, _ = pathhelper.GetCurrentExecDir()
	configFile = path.Join(currentDir, "config.json")
}

// loadConfig loads app config.
func loadConfig(configFile string, config *Config) error {
	// Load Conifg
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("load config file error: %v", err)

	}

	if err = json.Unmarshal(buf, config); err != nil {
		return fmt.Errorf("parse config err: %v", err)
	}

	return nil
}

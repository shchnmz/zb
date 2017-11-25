package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"

	"github.com/northbright/pathhelper"
	"github.com/shchnmz/zb"
)

// Config represents app settings.
type Config struct {
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

var (
	config Config
)

func main() {
	var (
		err error
	)

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	if err = loadConfig("config.json", &config); err != nil {
		return
	}

	if err = zb.ClearBlacklist(config.RedisServer, config.RedisPassword); err != nil {
		return
	}
}

func loadConfig(file string, config *Config) error {
	var (
		err        error
		buf        []byte
		currentDir string
	)

	currentDir, _ = pathhelper.GetCurrentExecDir()
	file = path.Join(currentDir, file)

	// Load Conifg
	if buf, err = ioutil.ReadFile(file); err != nil {
		return err
	}

	if err = json.Unmarshal(buf, &config); err != nil {
		return err
	}

	return nil
}

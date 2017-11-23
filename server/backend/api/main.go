package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/northbright/pathhelper"
)

type Config struct {
	ServerAddr    string `json:"server_addr"`
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

var (
	currentDir, configFile string
	config                 Config
	blacklists             map[string][]string
)

func main() {
	var (
		err error
	)

	defer func() {
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}()

	if err = loadConfig(); err != nil {
		err = fmt.Errorf("loadConfig() error: %v", err)
		return
	}

	// Load blacklists.
	if blacklists, err = loadBlacklists(); err != nil {
		err = fmt.Errorf("loadBacklists() error: %v", err)
		return
	}

	for k, v := range blacklists {
		fmt.Printf("blacklist: %v\n", k)
		for _, data := range v {
			fmt.Printf("%v\n", data)
		}
	}

	r := gin.Default()

	// Core APIs.

	//r.POST("/employee", addEmployee)
	// Set employee
	//r.PUT("/employee/:id", setEmployee)

	// Get student names by phone num.
	r.GET("/get-names-by-phone-num/:phone_num", getNamesByPhoneNum)

	// Get classes by name and phone num.
	r.GET("/get-classes-by-name-and-phone-num/:name/:phone_num", getClassesByNameAndPhoneNum)

	r.Run(config.ServerAddr)

}

// init initializes path variables.
func init() {
	currentDir, _ = pathhelper.GetCurrentExecDir()
	configFile = path.Join(currentDir, "config.json")
}

// loadConfig loads app config.
func loadConfig() error {
	// Load Conifg
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("load config file error: %v", err)

	}

	if err = json.Unmarshal(buf, &config); err != nil {
		return fmt.Errorf("parse config err: %v", err)
	}

	return nil
}

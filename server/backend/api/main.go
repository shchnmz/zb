package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/northbright/pathhelper"
	"github.com/shchnmz/ming"
	"github.com/shchnmz/zb"
)

type Config struct {
	ServerAddr    string `json:"server_addr"`
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

var (
	currentDir, configFile string
	config                 Config
	db                     zb.DB
)

func main() {
	var (
		err error
	)

	defer func() {
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}()

	if err = loadConfig(); err != nil {
		err = fmt.Errorf("loadConfig() error: %v", err)
		return
	}

	// Init DB.
	db = zb.DB{ming.DB{config.RedisServer, config.RedisPassword}}

	r := gin.Default()

	// Core APIs.

	// Get student names by phone num.
	r.GET("/get-names-by-phone-num/:phone_num", getNamesByPhoneNum)

	// Get classes by name and phone num.
	r.GET("/get-classes-by-name-and-phone-num/:name/:phone_num", getClassesByNameAndPhoneNum)

	// Get available periods for the category of the class.
	r.GET("/get-available-periods/:class", getAvailablePeriods)

	// Post request.
	r.POST("/request", postRequest)

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

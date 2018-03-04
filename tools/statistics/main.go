package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/northbright/maphelper"
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
		err error
	)

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	if err = loadConfig(configFile, &config); err != nil {
		return
	}

	if err = getStatistics(config.RedisServer, config.RedisPassword); err != nil {
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

func getStatistics(redisServer, redisPassword string) error {
	db := &zb.DB{ming.DB{redisServer, redisPassword}}
	s, err := db.GetStatistics()
	if err != nil {
		return err
	}

	fmt.Printf("按校区统计:\n")
	keys, err := maphelper.SortMapByValues(s.StudentNumOfEachCampus, true)
	if err != nil {
		return err
	}

	for _, key := range keys {
		fmt.Printf("%v: %v人\n", key, s.StudentNumOfEachCampus[key])
	}

	fmt.Printf("\n按年级统计:\n")
	keys, err = maphelper.SortMapByValues(s.StudentNumOfEachCategory, true)
	if err != nil {
		return err
	}

	for _, key := range keys {
		fmt.Printf("%v: %v人\n", key, s.StudentNumOfEachCategory[key])
	}

	fmt.Printf("\n按教师统计:\n")
	keys, err = maphelper.SortMapByValues(s.StudentNumOfEachTeacher, true)
	if err != nil {
		return err
	}

	for _, key := range keys {
		fmt.Printf("%v: %v人\n", key, s.StudentNumOfEachTeacher[key])
	}

	fmt.Printf("\n按教师转班率统计:\n")
	keys, err = maphelper.SortMapByValues(s.StudentPercentOfEachTeacher, true)
	if err != nil {
		return err
	}

	for _, key := range keys {
		fmt.Printf("%v: %.2f%%\n", key, s.StudentPercentOfEachTeacher[key])
	}

	return nil
}

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

	if err = getStatistics(config.RedisServer, config.RedisPassword); err != nil {
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

	return json.Unmarshal(buf, &config)
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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/garyburd/redigo/redis"
	"github.com/northbright/pathhelper"
	"github.com/northbright/redishelper"
)

type Config struct {
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

var (
	config Config
)

func main() {
	var (
		err                    error
		buf                    []byte
		currentDir, configFile string
	)

	defer func() {
		if err != nil {
			fmt.Printf("%v", err)
		}
	}()

	currentDir, _ = pathhelper.GetCurrentExecDir()
	configFile = path.Join(currentDir, "config.json")

	// Load Conifg
	if buf, err = ioutil.ReadFile(configFile); err != nil {
		err = fmt.Errorf("load config file error: %v", err)
		return
	}

	if err = json.Unmarshal(buf, &config); err != nil {
		err = fmt.Errorf("parse config err: %v", err)
		return
	}

	err = listAllClasses(config.RedisServer, config.RedisPassword)
}

// listAllClasses lists all classes in ming800.
func listAllClasses(redisServer, redisPassword string) error {
	var (
		err error
	)

	conn, err := redishelper.GetRedisConn(redisServer, redisPassword)
	if err != nil {
		return err
	}
	defer conn.Close()

	k := "campuses"
	campuses, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return err
	}

	for _, campus := range campuses {
		k = fmt.Sprintf("%v:categories", campus)
		categories, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
		if err != nil {
			return err
		}

		for _, category := range categories {
			k = fmt.Sprintf("%v:%v:periods", campus, category)
			periods, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
			if err != nil {
				return err
			}

			for _, period := range periods {
				fmt.Printf("%v,%v,%v\n", campus, category, period)
			}

		}
	}

	return nil
}

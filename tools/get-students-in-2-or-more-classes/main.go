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

	err = FindStudents(config.RedisServer, config.RedisPassword)
}

// FindStudents find the students which are in 2 or more classes then output student's name, phone and classes.
func FindStudents(redisServer, redisPassword string) error {
	var (
		err   error
		v     []interface{}
		items []string
	)

	conn, err := redishelper.GetRedisConn(redisServer, redisPassword)
	if err != nil {
		return err
	}
	defer conn.Close()

	k := "students"
	cursor := 0
	for {
		if v, err = redis.Values(conn.Do("ZSCAN", k, cursor, "COUNT", 1000)); err != nil {
			return err
		}

		if _, err = redis.Scan(v, &cursor, &items); err != nil {
			return err
		}

		l := len(items)
		if l <= 0 || l%2 != 0 {
			continue
		}

		for i := 0; i < l; i += 2 {
			key := fmt.Sprintf("%v:classes", items[i])
			count, err := redis.Int64(conn.Do("ZCARD", key))
			if err != nil {
				return err
			}

			if count < 2 {
				continue
			}

			// Output student name / phone num.
			fmt.Printf("%v\n", items[i])
			classes, err := redis.Strings(conn.Do("ZRANGE", key, 0, -1))
			if err != nil {
				return err
			}

			// Output student's  classes.
			for _, class := range classes {
				fmt.Printf("%v\n", class)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

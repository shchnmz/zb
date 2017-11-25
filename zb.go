package zb

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/northbright/redishelper"
)

// Config represents the settings.
type Config struct {
	RedisServer   string
	RedisPassword string
}

// ZB represents the transfer utility.
type ZB struct {
	Config
}

// NewZB creates a new instance of ZB.
func NewZB(redisServer, redisPassword string) ZB {
	return ZB{Config{redisServer, redisPassword}}
}

// GetNamesByPhoneNum searchs student names by phone number.
func (z *ZB) GetNamesByPhoneNum(phoneNum string) ([]string, error) {
	var (
		err   error
		names []string
	)

	conn, err := redishelper.GetRedisConn(z.RedisServer, z.RedisPassword)
	if err != nil {
		return names, err
	}
	defer conn.Close()

	k := fmt.Sprintf("ming:%v:students", phoneNum)
	names, err = redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return names, err
	}

	return names, nil
}

// GetClassesByNameAndPhoneNum searchs classes by student name and phone number.
func (z *ZB) GetClassesByNameAndPhoneNum(name, phoneNum string) ([]string, error) {
	var (
		err     error
		classes []string
	)

	conn, err := redishelper.GetRedisConn(z.RedisServer, z.RedisPassword)
	if err != nil {
		return classes, err
	}
	defer conn.Close()

	k := fmt.Sprintf("ming:%v:%v:classes", name, phoneNum)
	classes, err = redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return classes, err
	}

	return classes, nil
}

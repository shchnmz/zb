package zb

import (
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/northbright/redishelper"
)

// ValidClassString validates the class string.
//
// Params:
//     classStr: the class string should be: "$CAMPUS:$CATEGORY:$CLASS".
//               e.g. "新校区:一年级:一年级3班".
func (z *ZB) ValidClassString(classStr string) (bool, error) {
	conn, err := redishelper.GetRedisConn(z.RedisServer, z.RedisPassword)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	arr := strings.SplitN(classStr, ":", 3)
	campus := arr[0]
	category := arr[1]
	class := arr[2]

	k := fmt.Sprintf("ming:%v:%v:classes", campus, category)
	score, err := redis.String(conn.Do("ZSCORE", k, class))
	if err != nil {
		return false, err
	}

	if score == "" {
		return false, nil
	}

	return true, nil
}

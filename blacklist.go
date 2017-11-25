package zb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"time"

	"github.com/northbright/pathhelper"
	"github.com/northbright/redishelper"
)

// Blacklist represents the blacklist of student transter.
type Blacklist struct {
	List map[string][]string `json:"blacklist"`
}

var (
	blacklistTypes = map[string]string{
		"from_campuses": "can't transfer students from the campuses",
		"from_periods":  "can't transfer students from the periods",
		"from_classes":  "can't transfer students from classes",
		"to_campuses":   "can't transfer students to the campuses",
		"to_periods":    "can't transfer students to the periods",
		"to_classes":    "can't transfer students to the classes",
	}
)

// SetBlacklist updates the backlist in redis.
//
// Params:
//     redisServer: redis server address. e.g. ":6379"
//     redisPassword: redis password.
//     blacklist:
//       There're following types of blacklist:
//       "from_campuses", "from_periods", "from_classes",
//       "to_campuses", "to_periods", "to_classes".
func SetBlacklist(redisServer, redisPassword string, blacklist map[string][]string) error {
	pipedConn, err := redishelper.GetRedisConn(redisServer, redisPassword)
	if err != nil {
		return err
	}
	defer pipedConn.Close()

	if !ValidBlacklist(blacklist) {
		return fmt.Errorf("invalid blacklist")
	}

	pipedConn.Send("MULTI")

	for k, list := range blacklist {
		key := fmt.Sprintf("zb:blacklist:%v", k)

		// Delete key before update the list.
		pipedConn.Send("DEL", key)
		for _, data := range list {
			// Get timestamp as score for redis ordered set.
			t := strconv.FormatInt(time.Now().UnixNano(), 10)
			pipedConn.Send("ZADD", key, t, data)
		}
	}

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return err
	}

	return nil
}

// ClearBlacklist clear the backlist in redis.
//
// Params:
//     redisServer: redis server address. e.g. ":6379"
//     redisPassword: redis password.
func ClearBlacklist(redisServer, redisPassword string) error {
	pipedConn, err := redishelper.GetRedisConn(redisServer, redisPassword)
	if err != nil {
		return err
	}
	defer pipedConn.Close()

	pipedConn.Send("MULTI")

	for k, _ := range blacklistTypes {
		key := fmt.Sprintf("zb:blacklist:%v", k)
		pipedConn.Send("DEL", key)
	}

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return err
	}

	return nil
}

// ValidBlacklist validates the backlist.
func ValidBlacklist(blacklist map[string][]string) bool {
	for k, _ := range blacklist {
		if _, ok := blacklistTypes[k]; !ok {
			return false
		}
	}
	return true
}

// LoadBlacklist loads the blacklist from JSON file then set it to redis.
//
// Params:
//     redisServer: redis server address. e.g. ":6379".
//     redisPassword: redis password.
//     file: JSON file name.
//       There're following types of blacklist:
//       "from_campuses", "from_periods", "from_classes",
//       "to_campuses", "to_periods", "to_classes".
//       Example blacklist.json:
//       {
//         "blacklist": {
//           "from_campuses":["新校区"],
//           "from_periods":[],
//           "from_classes":[
//             "新校区:二年级:17秋新基二三1",
//             "新校区:二年级:17秋新基二三2"
//           ],
//           "to_campuses":["新校区"],
//           "to_periods": [
//             "老校区:幼中:星期二16:25-17:55",
//             "老校区:幼中:星期三16:25-17:55",,
//             "to_classes":[]
//         }
//       }
func LoadBlacklist(redisServer, redisPassword, file string, blacklist *Blacklist) error {
	var (
		err        error
		buf        []byte
		currentDir string
	)

	currentDir, _ = pathhelper.GetCurrentExecDir()
	file = path.Join(currentDir, file)

	// Load blacklist.
	if buf, err = ioutil.ReadFile(file); err != nil {
		return err
	}

	if err = json.Unmarshal(buf, blacklist); err != nil {
		return err
	}

	// Set blacklist to redis.
	return SetBlacklist(redisServer, redisPassword, blacklist.List)
}

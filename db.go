package zb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/northbright/pathhelper"
	"github.com/northbright/redishelper"
	"github.com/shchnmz/ming"
)

// DB represents database to store transfer data.
// It's wrapper of ming.DB.
// Usage:
// import "github.com/shchnmz/ming"
// db := zb.DB{ming.DB{redisServer, redisPassword}}
type DB struct {
	ming.DB
}

// Blacklist represents the blacklist of student transter.
type Blacklist struct {
	List map[string][]string `json:"blacklist"`
}

// Record represents the transfer record.
type Record struct {
	Name       string `redis:"name"`
	PhoneNum   string `redis:"phone_num"`
	Category   string `redis:"category"`
	FromCampus string `redis:"from_campus"`
	FromClass  string `redis:"from_class"`
	FromPeriod string `redis:"from_period"`
	ToCampus   string `redis:"to_campus"`
	ToPeriod   string `redis:"to_period"`
	Time       string `redis:"time"`
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
//     blacklist:
//       There're following types of blacklist:
//       "from_campuses", "from_periods", "from_classes",
//       "to_campuses", "to_periods", "to_classes".
func (db *DB) SetBlacklist(blacklist map[string][]string) error {
	pipedConn, err := redishelper.GetRedisConn(db.RedisServer, db.RedisPassword)
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
func (db *DB) ClearBlacklist() error {
	pipedConn, err := redishelper.GetRedisConn(db.RedisServer, db.RedisPassword)
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
func (db *DB) LoadBlacklist(file string, blacklist *Blacklist) error {
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
	return db.SetBlacklist(blacklist.List)
}

// IsFromClassInBlacklist checks if the class transfer from is in blacklist.
func (db *DB) IsFromClassInBlacklist(campus, category, class string) (bool, error) {
	conn, err := redishelper.GetRedisConn(db.RedisServer, db.RedisPassword)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	// Workflow: check campus -> check period -> check class.
	// Step 1. Check if campus is in blacklist.
	k := "zb:blacklist:from_campuses"
	m := campus
	score, err := redis.String(conn.Do("ZSCORE", k, m))
	if err != nil && err != redis.ErrNil {
		return false, err
	}

	if score != "" {
		return true, nil
	}

	// Get period if the class.
	period, err := db.GetClassPeriod(campus, category, class)
	if err != nil {
		return false, err
	}

	// Step 2. Check if period is in blacklist.
	k = "zb:blacklist:from_periods"
	m = fmt.Sprintf("%v:%v:%v", campus, category, period)
	score, err = redis.String(conn.Do("ZSCORE", k, m))
	if err != nil && err != redis.ErrNil {
		return false, err
	}

	if score != "" {
		return true, nil
	}

	// Step 3. Check if class is in blacklist.
	k = "zb:blacklist:from_classes"
	m = fmt.Sprintf("%v:%v:%v", campus, category, class)
	score, err = redis.String(conn.Do("ZSCORE", k, m))
	if err != nil && err != redis.ErrNil {
		return false, err
	}

	if score != "" {
		return true, nil
	}

	return false, nil
}

// IsToPeriodInBlacklist checks if the period transfer to is in blacklist.
func (db *DB) IsToPeriodInBlacklist(campus, category, period string) (bool, error) {
	conn, err := redishelper.GetRedisConn(db.RedisServer, db.RedisPassword)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	// Workflow: check campus -> check period.
	// Step 1. Check if campus is in blacklist.
	k := "zb:blacklist:to_campuses"
	m := campus
	score, err := redis.String(conn.Do("ZSCORE", k, m))
	if err != nil && err != redis.ErrNil {
		return false, err
	}

	if score != "" {
		return true, nil
	}

	// Step 2. Check if period is in blacklist.
	k = "zb:blacklist:to_periods"
	m = fmt.Sprintf("%v:%v:%v", campus, category, period)
	score, err = redis.String(conn.Do("ZSCORE", k, m))
	if err != nil && err != redis.ErrNil {
		return false, err
	}

	if score != "" {
		return true, nil
	}

	return false, nil
}

// FilterToPeriodsOfCategoryWithBlacklist filters periods of the category in blacklist.
//
// Params:
//     category: category of the periods.
//     campusPeriods: map contains periods to filter. key: campus, value: periods.
// Returns:
//     filtered campus - periods map. key: campus, value: periods.
func (db *DB) FilterToPeriodsOfCategoryWithBlacklist(category string, campusPeriods map[string][]string) (map[string][]string, error) {
	conn, err := redishelper.GetRedisConn(db.RedisServer, db.RedisPassword)
	if err != nil {
		return map[string][]string{}, err
	}
	defer conn.Close()

	filteredCampusPeriods := map[string][]string{}
	for campus, periods := range campusPeriods {
		// Step 1. Check if campus is in blacklist.
		k := "zb:blacklist:to_campuses"
		m := campus
		score, err := redis.String(conn.Do("ZSCORE", k, m))
		if err != nil && err != redis.ErrNil {
			return map[string][]string{}, err
		}

		// Skip if campus is in blacklist.
		if score != "" {
			continue
		}

		// Step 2. Check if period is in blacklist.
		for _, period := range periods {
			k := "zb:blacklist:to_periods"
			m := fmt.Sprintf("%v:%v:%v", campus, category, period)
			score, err := redis.String(conn.Do("ZSCORE", k, m))
			if err != nil && err != redis.ErrNil {
				return map[string][]string{}, err
			}

			if score == "" {
				filteredCampusPeriods[campus] = append(filteredCampusPeriods[campus], period)
			}
		}
	}
	return filteredCampusPeriods, nil
}

// GetAvailblePeriodsOfCategory gets category's periods for all campuses, filtered with blacklist.
//
// Params:
//     category: category which you want to get all periods.
// Returns:
//     a map contains all periods. key: campus, value: periods.
func (db *DB) GetAvailblePeriodsOfCategory(category string) (map[string][]string, error) {
	// Get all periods of the category.
	campusPeriods, err := db.GetAllPeriodsOfCategory(category)
	if err != nil {
		return map[string][]string{}, err
	}

	filteredCampusPeriods, err := db.FilterToPeriodsOfCategoryWithBlacklist(category, campusPeriods)
	if err != nil {
		return map[string][]string{}, err
	}

	return filteredCampusPeriods, nil
}

// SetRecord sets the record in redis.
func (db *DB) SetRecord(r Record) error {
	pipedConn, err := redishelper.GetRedisConn(db.RedisServer, db.RedisPassword)
	if err != nil {
		return err
	}
	defer pipedConn.Close()

	pipedConn.Send("MULTI")

	k := fmt.Sprintf("zb:record:%v:%v", r.Name, r.PhoneNum)
	pipedConn.Send("HMSET", k, "name", r.Name, "phone_num", r.PhoneNum, "category", r.Category, "from_campus", r.FromCampus, "from_class", r.FromClass, "from_period", r.FromPeriod, "to_campus", r.ToCampus, "to_period", r.ToPeriod, "time", r.Time)

	timestamp := time.Now().Unix()
	key := "zb:records"
	m := k
	pipedConn.Send("ZADD", key, timestamp, m)

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return err
	}
	return nil
}

func (db *DB) GetAllRecords() ([]Record, error) {
	var records []Record

	conn, err := redishelper.GetRedisConn(db.RedisServer, db.RedisPassword)
	if err != nil {
		return []Record{}, err
	}
	defer conn.Close()

	k := "zb:records"
	keys, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return []Record{}, err
	}

	log.Printf("keys: %v", keys)
	for _, key := range keys {
		values, err := redis.Values(conn.Do("HGETALL", key))
		if err != nil {
			log.Printf("values error: %v", err)
			return []Record{}, err
		}

		record := Record{}
		if err = redis.ScanStruct(values, &record); err != nil {
			log.Printf("Scanstruct() error: %v", err)
			return []Record{}, err
		}

		records = append(records, record)
	}
	return records, nil

}

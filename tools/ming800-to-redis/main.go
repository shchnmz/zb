package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"github.com/garyburd/redigo/redis"
	"github.com/northbright/ming800"
	"github.com/northbright/pathhelper"
	"github.com/northbright/redishelper"
)

type Config struct {
	ServerURL     string `json:"server_url"`
	Company       string `json:"company"`
	User          string `json:"user"`
	Password      string `json:"password"`
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

var (
	config Config
)

// Valid8DigitTelephoneNum checks if phone number matches the format:
// 1. Starts with 8 digital number.
// 2. Can have one or more '.' as sufix.
func Valid8DigitTelephoneNum(phoneNum string) bool {
	p := `^\d{8}\.*$`
	re := regexp.MustCompile(p)
	return re.MatchString(phoneNum)
}

// ValidMobilePhoneNum checks if phone number matches the format:
// 1. Starts with 11 digital number.
// 2. Can have one or more '.' as sufix.
func ValidMobilePhoneNum(phoneNum string) bool {
	p := `^\d{11}\.*$`
	re := regexp.MustCompile(p)
	return re.MatchString(phoneNum)
}

// ValidPhoneNum checks if phone number is 11-digit mobile phone number or 8-digit telephone number.
func ValidPhoneNum(phoneNum string) bool {
	if !Valid8DigitTelephoneNum(phoneNum) && !ValidMobilePhoneNum(phoneNum) {
		return false
	}
	return true
}

// ParseCategory gets campus and real category from category string.
//
//   Param:
//       category: raw category string like this: 初一（中山）
//   Return:
//       campus, category. e.g. campus: 中山,category: 初一
func ParseCategory(category string) (string, string) {
	p := `^(\S+)（(\S+)）$`
	re := regexp.MustCompile(p)
	matched := re.FindStringSubmatch(category)
	if len(matched) != 3 {
		return "", ""
	}
	return matched[2], matched[1]
}

func main() {
	var (
		err                    error
		buf                    []byte
		currentDir, configFile string
		s                      *ming800.Session
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

	// New a session
	if s, err = ming800.NewSession(config.ServerURL, config.Company, config.User, config.Password); err != nil {
		err = fmt.Errorf("NewSession() error: %v", err)
		return
	}

	// Login
	if err = s.Login(); err != nil {
		err = fmt.Errorf("Login() error: %v", err)
		return
	}

	// Walk
	// Write your own class and student handler functions.
	// Class and student handler will be called while walking ming800.
	if err = s.Walk(classHandler, studentHandler); err != nil {
		err = fmt.Errorf("Walk() error: %v", err)
		return
	}

	// Logout
	if err = s.Logout(); err != nil {
		err = fmt.Errorf("Logout() error: %v", err)
		return
	}
}

func classHandler(class ming800.Class) {
	var err error

	defer func() {
		if err != nil {
			log.Printf("classHandler() error: %v", err)
		}
	}()

	pipedConn, err := redishelper.GetRedisConn(config.RedisServer, config.RedisPassword)
	if err != nil {
		return
	}
	defer pipedConn.Close()

	pipedConn.Do("MULTI")

	campus, category := ParseCategory(class.Category)
	if category == "" && campus == "" {
		err = fmt.Errorf("Failed to parse category and campus: %v", class.Category)
		return
	}

	// Get timestamp as score for redis ordered set.
	t := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Update SET: key: campus, value: categories.
	k := fmt.Sprintf("%v:categories", campus)
	pipedConn.Send("ZADD", k, t, category)

	// Update SET: key: category, value: campuses.
	k = fmt.Sprintf("%v:campuses", category)
	pipedConn.Send("ZADD", k, t, campus)

	// Update SET: key: campus + category, value: classes.
	k = fmt.Sprintf("%v:%v:classes", campus, category)
	pipedConn.Send("ZADD", k, t, class.Name)

	// Update SET: key: campus + category + class, value: teachers.
	k = fmt.Sprintf("%v:%v:%v:teachers", campus, category, class.Name)
	for _, teacher := range class.Teachers {
		t = strconv.FormatInt(time.Now().UnixNano(), 10)
		pipedConn.Send("ZADD", k, t, teacher)
	}

	// Update SET: key: campus + category + class, value: periods.
	k = fmt.Sprintf("%v:%v:%v:periods", campus, category, class.Name)
	for _, period := range class.Periods {
		t = strconv.FormatInt(time.Now().UnixNano(), 10)
		pipedConn.Send("ZADD", k, t, period)
	}

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return
	}
}

func studentHandler(class ming800.Class, student ming800.Student) {
	var err error

	defer func() {
		if err != nil {
			log.Printf("studentHandler() error: %v", err)
		}
	}()

	// Check if phone number: 11-digit or 8-digit.
	if !ValidPhoneNum(student.PhoneNum) {
		fmt.Printf("%s,%s,%s,%s\n", class.Category, class.Name, student.Name, student.PhoneNum)
		return
	}

	// Student contact phone may have '.' suffix, remove it.
	student.PhoneNum = strings.TrimRight(student.PhoneNum, `.`)

	// Get another redis connection for pipelined transaction.
	pipedConn, err := redishelper.GetRedisConn(config.RedisServer, config.RedisPassword)
	if err != nil {
		return
	}
	defer pipedConn.Close()

	pipedConn.Do("MULTI")

	// Get timestamp as store for redis ordered set.
	t := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Update SET: key: student name + student contact phone num, value: classes.
	k := fmt.Sprintf("%v:%v:classes", student.Name, student.PhoneNum)
	pipedConn.Send("ZADD", k, t, class.Name)

	// Update SET: key: student name, value: student contact phone numbers.
	k = fmt.Sprintf("%v:phones", student.Name)
	pipedConn.Send("ZADD", k, t, student.PhoneNum)

	// Update SET: key: student contact phone num, value: student names.
	k = fmt.Sprintf("%v:names", student.PhoneNum)
	pipedConn.Send("ZADD", k, t, student.Name)

	campus, category := ParseCategory(class.Category)
	if category == "" && campus == "" {
		err = fmt.Errorf("Failed to parse category and campus: %v", class.Category)
		return
	}

	// Update SET: key: campus + category + class, value: student name + student contact phone num.
	k = fmt.Sprintf("%v:%v:%v:students", campus, category, class.Name)
	v := fmt.Sprintf("%v:%v", student.Name, student.PhoneNum)
	pipedConn.Send("ZADD", k, t, v)

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return
	}
}

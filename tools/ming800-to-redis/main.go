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

	// Add category name to redis set.
	k := "categories"
	t := strconv.FormatInt(time.Now().UnixNano(), 10)
	pipedConn.Send("ZADD", k, t, class.Category)

	k = fmt.Sprintf("category:%v:classes", class.Category)
	pipedConn.Send("ZADD", k, t, class.Name)

	k = fmt.Sprintf("class:%v:%v:teachers", class.Category, class.Name)
	for _, teacher := range class.Teachers {
		t = strconv.FormatInt(time.Now().UnixNano(), 10)
		pipedConn.Send("ZADD", k, t, teacher)
	}

	k = fmt.Sprintf("class:%v:%v:periods", class.Category, class.Name)
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

	k := fmt.Sprintf("student:%v:%v:classes", student.Name, student.PhoneNum)
	t := strconv.FormatInt(time.Now().UnixNano(), 10)
	pipedConn.Send("ZADD", k, t, class.Name)

	k = fmt.Sprintf("student:%v:phones", student.Name)
	t = strconv.FormatInt(time.Now().UnixNano(), 10)
	pipedConn.Send("ZADD", k, t, student.PhoneNum)

	k = fmt.Sprintf("student:%v:names", student.PhoneNum)
	t = strconv.FormatInt(time.Now().UnixNano(), 10)
	pipedConn.Send("ZADD", k, t, student.Name)

	k = fmt.Sprintf("class:%v:%v:students", class.Category, class.Name)
	t = strconv.FormatInt(time.Now().UnixNano(), 10)
	v := fmt.Sprintf("%v:%v", student.Name, student.PhoneNum)
	pipedConn.Send("ZADD", k, t, v)

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return
	}
}

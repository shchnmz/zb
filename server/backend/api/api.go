package main

import (
	"log"
	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/shchnmz/ming"
	"github.com/shchnmz/zb"
)

// TransferRequest represents the transfer request information.
type TransferRequest struct {
	Name       string `json:"name" binding:"required"`
	PhoneNum   string `json:"phone_num" binding:"required"`
	FromClass  string `json:"from_class" binding:"required"`
	ToCampus   string `json:"to_campus" binding:"required"`
	ToCategory string `json:"to_category" binding:"required"`
	ToPeriod   string `json:"to_period" binding:"required"`
}

func getNamesByPhoneNum(c *gin.Context) {
	var (
		err        error
		statusCode = 200
		errMsg     = ""
		names      []string
	)

	defer func() {
		if err != nil {
			log.Printf("getNamesByPhone() err: %v", err)
		}

		if errMsg != "" {
			log.Printf("getNamesByPhone() errMsg: %v", errMsg)
		}

		c.JSON(statusCode, gin.H{"err_msg": errMsg, "data": gin.H{"names": names}})
	}()

	phoneNum := c.Param("phone_num")
	if !ming.ValidPhoneNum(phoneNum) {
		statusCode = 400
		errMsg = "invalid phone num"
		return
	}

	z := zb.NewZB(config.RedisServer, config.RedisPassword)
	if names, err = z.GetNamesByPhoneNum(phoneNum); err != nil {
		statusCode = 500
		errMsg = "internal server error"
		return
	}

	log.Printf("phone num: %v, names: %v", phoneNum, names)
}

func getClassesByNameAndPhoneNum(c *gin.Context) {
	var (
		err        error
		statusCode = 200
		errMsg     = ""
		classes    []string
	)

	defer func() {
		if err != nil {
			log.Printf("getClassByNameAndPhoneNum() err: %v", err)
		}

		if errMsg != "" {
			log.Printf("getClassByNameAndPhoneNum() errMsg: %v", errMsg)
		}

		c.JSON(statusCode, gin.H{"err_msg": errMsg, "data": gin.H{"classes": classes}})
	}()

	name := c.Param("name")
	if len(name) == 0 {
		statusCode = 400
		errMsg = "empty name"
		return
	}

	phoneNum := c.Param("phone_num")
	if !ming.ValidPhoneNum(phoneNum) {
		statusCode = 400
		errMsg = "Invalid phone num"
		return
	}

	z := zb.NewZB(config.RedisServer, config.RedisPassword)
	if classes, err = z.GetClassesByNameAndPhoneNum(name, phoneNum); err != nil {
		statusCode = 500
		errMsg = "internal server error"
		return
	}

	log.Printf("name: %v, phone num: %v, classes: %v", name, phoneNum, classes)
}

/*
func getAvailablePeriodsByClass(c *gin.Context) {
	var (
		err        error
		statusCode = 200
		errMsg     = ""
		// campusPeriods is the campus -> periods map.
		campusPeriods = map[string][]string{}
	)

	defer func() {
		if err != nil {
			log.Printf("getClassByNameAndPhoneNum() err: %v", err)
		}

		if errMsg != "" {
			log.Printf("getClassByNameAndPhoneNum() errMsg: %v", errMsg)
		}

		c.JSON(statusCode, gin.H{"err_msg": errMsg, "data": gin.H{"campus_periods": campusPeriods}})
	}()

	class := c.Param("class")
	if len(class) == 0 {
		err = fmt.Errorf("empty class")
		statusCode = 400
		errMsg = "empty name"
		return
	}

	arr := strings.SplitN(class, ":", 3)
	campus := arr[0]
	category := arr[1]
}
*/

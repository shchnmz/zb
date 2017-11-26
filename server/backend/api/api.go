package main

import (
	"log"
	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/shchnmz/ming"
	//	"github.com/shchnmz/zb"
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

	if names, err = db.GetNamesByPhoneNum(phoneNum); err != nil {
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

	if classes, err = db.GetClassesByNameAndPhoneNum(name, phoneNum); err != nil {
		statusCode = 500
		errMsg = "internal server error"
		return
	}

	log.Printf("name: %v, phone num: %v, classes: %v", name, phoneNum, classes)
}

func getAvailablePeriods(c *gin.Context) {
	var (
		err        error
		statusCode = 200
		errMsg     = ""
		// campusPeriods is the campus -> periods map.
		campusPeriods = map[string][]string{}
	)

	defer func() {
		if err != nil {
			log.Printf("getAvailablePeriods() err: %v", err)
		}

		if errMsg != "" {
			log.Printf("getAvailablePeriods() errMsg: %v", errMsg)
		}

		c.JSON(statusCode, gin.H{"err_msg": errMsg, "data": gin.H{"campus_periods": campusPeriods}})
	}()

	classValue := c.Param("class")
	if len(classValue) == 0 {
		statusCode = 400
		errMsg = "empty name"
		return
	}

	campus, category, class := ming.ParseClassValue(classValue)
	valid, err := db.ValidClass(campus, category, class)
	if err != nil {
		statusCode = 500
		errMsg = "internal server error"
		return
	}

	if !valid {
		statusCode = 400
		errMsg = "invalid class value"
		return
	}

	// Check if campus is in from_campuses blacklist.
	in, err := db.IsFromClassInBlacklist(campus, category, class)
	if err != nil {
		statusCode = 500
		errMsg = "internal server error"
		return
	}

	if in {
		statusCode = 400
		errMsg = "当前班级不在本次转班范围内"
		return
	}

	if campusPeriods, err = db.GetAllPeriodsOfCategory(category); err != nil {
		statusCode = 500
		errMsg = "internal server error"
		return
	}

	log.Printf("class value: %v, campusPeriods: %v", classValue, campusPeriods)
}

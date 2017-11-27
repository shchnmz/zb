package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shchnmz/ming"
	"github.com/shchnmz/zb"
)

// Request represents the transfer request information.
type Request struct {
	Name      string `json:"name" binding:"required"`
	PhoneNum  string `json:"phone_num" binding:"required"`
	FromClass string `json:"from_class" binding:"required"`
	ToCampus  string `json:"to_campus" binding:"required"`
	ToPeriod  string `json:"to_period" binding:"required"`
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
		errMsg = "无效联系电话"
		return
	}

	if names, err = db.GetNamesByPhoneNum(phoneNum); err != nil {
		statusCode = 500
		errMsg = "服务器错误"
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
		errMsg = "姓名为空"
		return
	}

	phoneNum := c.Param("phone_num")
	if !ming.ValidPhoneNum(phoneNum) {
		statusCode = 400
		errMsg = "无效联系电话"
		return
	}

	if classes, err = db.GetClassesByNameAndPhoneNum(name, phoneNum); err != nil {
		statusCode = 500
		errMsg = "服务器错误"
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
		errMsg = "转出班级为空"
		return
	}

	campus, category, class := ming.ParseClassValue(classValue)
	valid, err := db.ValidClass(campus, category, class)
	if err != nil {
		statusCode = 500
		errMsg = "服务器错误"
		return
	}

	if !valid {
		statusCode = 400
		errMsg = "无效的转出班级"
		return
	}

	// Check if campus is in from_campuses blacklist.
	in, err := db.IsFromClassInBlacklist(campus, category, class)
	if err != nil {
		statusCode = 500
		errMsg = "服务器错误"
		return
	}

	if in {
		statusCode = 400
		errMsg = "转出班级不在本次转班范围内"
		return
	}

	// Get all periods of the category.
	if campusPeriods, err = db.GetAvailblePeriodsOfCategory(category); err != nil {
		statusCode = 500
		errMsg = "服务器错误"
		return
	}

	log.Printf("class value: %v, campusPeriods: %v", classValue, campusPeriods)
}

func postRequest(c *gin.Context) {
	var (
		err        error
		statusCode = 200
		errMsg     = ""
		request    Request
		record     zb.Record
	)

	defer func() {
		if err != nil {
			log.Printf("postRequest() err: %v", err)
		}

		if errMsg != "" {
			log.Printf("postRequest() errMsg: %v", errMsg)
		}

		c.JSON(statusCode, gin.H{"err_msg": errMsg, "data": gin.H{"request": request, "record": record}})
	}()

	if err = c.BindJSON(&request); err != nil {
		statusCode = 400
		errMsg = "无效转班请求"
		return
	}

	// Validte student name and phone num(see if the class can be found by name and phone num).
	classes, err := db.GetClassesByNameAndPhoneNum(request.Name, request.PhoneNum)
	if err != nil {
		statusCode = 500
		errMsg = "服务器错误"
		return
	}

	if len(classes) == 0 {
		statusCode = 400
		errMsg = "无法找到联系电话和学生姓名对应的班级"
		return
	}

	// Check if from class is in the classes of this student.
	found := false
	for _, v := range classes {
		if v == request.FromClass {
			found = true
		}
	}

	if !found {
		statusCode = 400
		errMsg = "学生姓名手机与转出班级信息不符"
	}

	// Validate class.
	fromCampus, category, fromClass := ming.ParseClassValue(request.FromClass)
	valid, err := db.ValidClass(fromCampus, category, fromClass)
	if err != nil {
		statusCode = 500
		errMsg = "服务器错误"
		return
	}

	if !valid {
		statusCode = 400
		errMsg = "无效的转出班级"
		return
	}

	// Get period to transfer from.
	fromPeriod, err := db.GetClassPeriod(fromCampus, category, fromClass)
	if err != nil {
		statusCode = 500
		errMsg = "服务器错误"
		return
	}

	// Validate the period to transfer to.
	valid, err = db.ValidPeriod(request.ToCampus, category, request.ToPeriod)
	if err != nil {
		statusCode = 500
		errMsg = "服务器错误"
		return
	}

	if !valid {
		statusCode = 400
		errMsg = "无效的转入时间段"
		return
	}

	// Check if from campus / period is equal to to campus / period.
	if fromCampus == request.ToCampus && fromPeriod == request.ToPeriod {
		statusCode = 400
		errMsg = "转出时间段与转入时间段相同，无效请求"
		return
	}

	// Check if campus is in from_campuses blacklist.
	in, err := db.IsFromClassInBlacklist(fromCampus, category, fromClass)
	if err != nil {
		statusCode = 500
		errMsg = "服务器错误"
		return
	}

	if in {
		statusCode = 400
		errMsg = "转出班级不在本次转班范围内"
		return
	}

	in, err = db.IsToPeriodInBlacklist(request.ToCampus, category, request.ToPeriod)
	if err != nil {
		statusCode = 500
		errMsg = "服务器错误"
		return
	}

	if in {
		statusCode = 400
		errMsg = "转入班级不在本次转班范围内"
		return
	}

	t := time.Now().Local()
	tm := fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	record = zb.Record{request.Name, request.PhoneNum, category, fromCampus, fromClass, fromPeriod, request.ToCampus, request.ToPeriod, tm}
	if err = db.SetRecord(record); err != nil {
		statusCode = 500
		errMsg = "服务器错误"
		return
	}

	log.Printf("request: %v, record: %v", request, record)
}

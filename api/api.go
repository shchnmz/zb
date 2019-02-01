package main

import (
	"fmt"
	"log"
	"strings"
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

// getNamesByPhoneNum is the handler to get the student names by phone number.
func getNamesByPhoneNum(c *gin.Context) {
	var (
		err     error
		success = false
		errMsg  = ""
		names   []string
	)

	defer func() {
		if err != nil {
			errMsg = "服务器错误"
			log.Printf("getNamesByPhone() err: %v", err)
		}

		if errMsg != "" {
			log.Printf("getNamesByPhone() errMsg: %v", errMsg)
		}

		c.JSON(200, gin.H{"success": success, "err_msg": errMsg, "names": names})
	}()

	phoneNum := c.Param("phone_num")
	if !ming.ValidPhoneNum(phoneNum) {
		errMsg = "无效联系电话"
		return
	}

	if names, err = db.GetNamesByPhoneNum(phoneNum); err != nil {
		return
	}

	success = true
	log.Printf("phone num: %v, names: %v", phoneNum, names)
}

// getClassesByNameAndPhoneNum is the handler to get classes by name and phone number.
func getClassesByNameAndPhoneNum(c *gin.Context) {
	var (
		err     error
		success = false
		errMsg  = ""
		classes []string
	)

	defer func() {
		if err != nil {
			log.Printf("getClassByNameAndPhoneNum() err: %v", err)
		}

		if errMsg != "" {
			log.Printf("getClassByNameAndPhoneNum() errMsg: %v", errMsg)
		}

		c.JSON(200, gin.H{"success": success, "err_msg": errMsg, "classes": classes})
	}()

	name := c.Param("name")
	if len(name) == 0 {
		errMsg = "姓名为空"
		return
	}

	phoneNum := c.Param("phone_num")
	if !ming.ValidPhoneNum(phoneNum) {
		errMsg = "无效联系电话"
		return
	}

	if classes, err = db.GetClassesByNameAndPhoneNum(name, phoneNum); err != nil {
		return
	}

	success = true
	log.Printf("name: %v, phone num: %v, classes: %v", name, phoneNum, classes)
}

// getTeachersByClass is the handler to get teachers by class name.
func getTeachersByClass(c *gin.Context) {
	var (
		err      error
		success  = false
		errMsg   = ""
		teachers []string
	)

	defer func() {
		if err != nil {
			log.Printf("getTeachersByClass() err: %v", err)
		}

		if errMsg != "" {
			log.Printf("getTeachersByClass() errMsg: %v", errMsg)
		}

		c.JSON(200, gin.H{"success": success, "err_msg": errMsg, "teachers": teachers})
	}()

	classWithCampusAndCategory := c.Param("class")
	if len(classWithCampusAndCategory) == 0 {
		errMsg = "班级为空"
		return
	}

	arr := strings.Split(classWithCampusAndCategory, ":")
	if len(arr) != 3 {
		errMsg = "班级格式错误"
		return
	}
	campus := arr[0]
	category := arr[1]
	class := arr[2]

	if teachers, err = db.GetTeachersOfClass(campus, category, class); err != nil {
		return
	}

	success = true
	log.Printf("campus: %v, category: %v, class: %v, teachers: %v", campus, category, class, teachers)
}

// getAvailablePeriods is the handler to get available periods by class string.
func getAvailablePeriods(c *gin.Context) {
	var (
		err     error
		success = false
		errMsg  = ""
		// campusPeriods is the campus -> periods map.
		campusPeriods = map[string][]string{}
	)

	defer func() {
		if err != nil {
			errMsg = "服务器错误"
			log.Printf("getAvailablePeriods() err: %v", err)
		}

		if errMsg != "" {
			log.Printf("getAvailablePeriods() errMsg: %v", errMsg)
		}

		c.JSON(200, gin.H{"success": success, "err_msg": errMsg, "campus_periods": campusPeriods})
	}()

	classValue := c.Param("class")
	if len(classValue) == 0 {
		errMsg = "转出班级为空"
		return
	}

	campus, category, class := ming.ParseClassValue(classValue)
	valid, err := db.ValidClass(campus, category, class)
	if err != nil {
		return
	}

	if !valid {
		errMsg = "无效的转出班级"
		return
	}

	// Check if campus is in from_campuses blacklist.
	in, err := db.IsFromClassInBlacklist(campus, category, class)
	if err != nil {
		return
	}

	if in {
		errMsg = "转出班级不在本次转班范围内"
		return
	}

	// Get all periods of the category.
	if campusPeriods, err = db.GetAvailblePeriodsOfCategory(category); err != nil {
		return
	}

	success = true
	log.Printf("class value: %v, campusPeriods: %v", classValue, campusPeriods)
}

// postRequest is the handler of transfer request.
func postRequest(c *gin.Context) {
	var (
		err     error
		success = false
		errMsg  = ""
		request Request
		record  zb.Record
	)

	defer func() {
		if err != nil {
			errMsg = "服务器错误"
			log.Printf("postRequest() err: %v", err)
		}

		if errMsg != "" {
			log.Printf("postRequest() errMsg: %v", errMsg)
		}

		c.JSON(200, gin.H{"success": success, "err_msg": errMsg, "request": request, "record": record})
	}()

	if err = c.BindJSON(&request); err != nil {
		errMsg = "无效转班请求"
		return
	}

	// Validte student name and phone num(see if the class can be found by name and phone num).
	classes, err := db.GetClassesByNameAndPhoneNum(request.Name, request.PhoneNum)
	if err != nil {
		return
	}

	if len(classes) == 0 {
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
		errMsg = "学生姓名手机与转出班级信息不符"
	}

	// Validate class.
	fromCampus, category, fromClass := ming.ParseClassValue(request.FromClass)
	valid, err := db.ValidClass(fromCampus, category, fromClass)
	if err != nil {
		return
	}

	if !valid {
		errMsg = "无效的转出班级"
		return
	}

	// Get period to transfer from.
	fromPeriod, err := db.GetClassPeriod(fromCampus, category, fromClass)
	if err != nil {
		return
	}

	// Validate the period to transfer to.
	valid, err = db.ValidPeriod(request.ToCampus, category, request.ToPeriod)
	if err != nil {
		return
	}

	if !valid {
		errMsg = "无效的转入时间段"
		return
	}

	// Check if from campus / period is equal to to campus / period.
	if fromCampus == request.ToCampus && fromPeriod == request.ToPeriod {
		errMsg = "转出时间段与转入时间段相同，无效请求"
		return
	}

	// Check if campus is in from_campuses blacklist.
	in, err := db.IsFromClassInBlacklist(fromCampus, category, fromClass)
	if err != nil {
		return
	}

	if in {
		errMsg = "转出班级不在本次转班范围内"
		return
	}

	in, err = db.IsToPeriodInBlacklist(request.ToCampus, category, request.ToPeriod)
	if err != nil {
		return
	}

	if in {
		errMsg = "转入班级不在本次转班范围内"
		return
	}

	t := time.Now().Local()
	tm := fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	record = zb.Record{request.Name, request.PhoneNum, category, fromCampus, fromClass, fromPeriod, request.ToCampus, request.ToPeriod, tm}
	if err = db.SetRecord(record); err != nil {
		return
	}

	success = true
	log.Printf("request: %v, record: %v", request, record)
}

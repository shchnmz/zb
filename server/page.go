package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/northbright/maphelper"
	"github.com/shchnmz/ming"
	"github.com/shchnmz/zb"
)

func home(c *gin.Context) {
	db := &zb.DB{ming.DB{config.RedisServer, config.RedisPassword}}

	enabled, err := db.IsEnabled()
	if err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"title": "转班申请系统错误",
		})
		log.Printf("home() db.IsEnabled() error: %v", err)
		return
	}

	if !enabled {
		c.HTML(http.StatusOK, "closed.tmpl", gin.H{
			"title":   "转班已经截止",
			"notices": config.ClosedNotices,
		})
	} else {

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "转班申请",
		})
	}
}

func admin(c *gin.Context) {
	db := &zb.DB{ming.DB{config.RedisServer, config.RedisPassword}}

	records, err := db.GetAllRecords()
	if err != nil {
		c.HTML(http.StatusOK, "error.tmpl", gin.H{
			"title": "转班申请系统错误",
		})
		log.Printf("admin() db.GetAllRecord() error: %v", err)
		return
	}

	c.HTML(http.StatusOK, "admin.tmpl", gin.H{
		"title":   "转班申请",
		"count":   len(records),
		"records": records,
	})
}

func statistics(c *gin.Context) {
	var (
		err   error
		items []string
	)

	defer func() {
		if err != nil {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"title": "转班申请系统错误",
			})
			log.Printf("statistics() error: %v", err)
			return
		}

		c.HTML(http.StatusOK, "statistics.tmpl", gin.H{
			"title": "转班统计",
			"items": items,
		})
	}()

	db := &zb.DB{ming.DB{config.RedisServer, config.RedisPassword}}
	s, err := db.GetStatistics()
	if err != nil {
		return
	}

	items = append(items, "--------------------------")
	items = append(items, "按校区统计:")
	items = append(items, "--------------------------")

	keys, err := maphelper.SortMapByValues(s.StudentNumOfEachCampus, true)
	if err != nil {
		return
	}

	for _, key := range keys {
		items = append(items, fmt.Sprintf("%v: %v人\n", key, s.StudentNumOfEachCampus[key]))
	}

	items = append(items, "--------------------------")
	items = append(items, "按年级统计:")
	items = append(items, "--------------------------")

	keys, err = maphelper.SortMapByValues(s.StudentNumOfEachCategory, true)
	if err != nil {
		return
	}

	for _, key := range keys {
		items = append(items, fmt.Sprintf("%v: %v人\n", key, s.StudentNumOfEachCategory[key]))
	}

	items = append(items, "--------------------------")
	items = append(items, "按教师转班率统计:")
	items = append(items, "--------------------------")

	keys, err = maphelper.SortMapByValues(s.StudentPercentOfEachTeacher, true)
	if err != nil {
		return
	}

	for _, key := range keys {
		items = append(items, fmt.Sprintf("%v: %.2f%%\n", key, s.StudentPercentOfEachTeacher[key]))
	}

	items = append(items, "--------------------------")
	items = append(items, "按教师统计:")
	items = append(items, "--------------------------")

	keys, err = maphelper.SortMapByValues(s.StudentNumOfEachTeacher, true)
	if err != nil {
		return
	}

	for _, key := range keys {
		items = append(items, fmt.Sprintf("%v: %v人\n", key, s.StudentNumOfEachTeacher[key]))
	}
}

func enable(c *gin.Context) {
	var (
		err     error
		enabled bool
	)

	defer func() {
		if err != nil {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"title": "转班申请系统错误",
			})
			log.Printf("enable() error: %v", err)
			return
		}

		c.HTML(http.StatusOK, "enable.tmpl", gin.H{
			"title":   "允许转班",
			"enabled": enabled,
		})
	}()

	db := &zb.DB{ming.DB{config.RedisServer, config.RedisPassword}}
	if err = db.Enable(true); err != nil {
		return
	}

	if enabled, err = db.IsEnabled(); err != nil {
		return
	}
}

func disable(c *gin.Context) {
	var (
		err     error
		enabled bool
	)

	defer func() {
		if err != nil {
			c.HTML(http.StatusOK, "error.tmpl", gin.H{
				"title": "转班申请系统错误",
			})
			log.Printf("disable() error: %v", err)
			return
		}

		c.HTML(http.StatusOK, "enable.tmpl", gin.H{
			"title":   "关闭转班",
			"enabled": enabled,
		})
	}()

	db := &zb.DB{ming.DB{config.RedisServer, config.RedisPassword}}
	if err = db.Enable(false); err != nil {
		return
	}

	if enabled, err = db.IsEnabled(); err != nil {
		return
	}
}

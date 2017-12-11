package main

import (
	//"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

package main

import (
	//"fmt"
	//"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shchnmz/ming"
	"github.com/shchnmz/zb"
)

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "转班申请",
	})
}

func admin(c *gin.Context) {
	var records []zb.Record

	defer func() {
		c.HTML(http.StatusOK, "admin.tmpl", gin.H{
			"title":   "转班申请",
			"count":   len(records),
			"records": records,
		})
	}()

	db := &zb.DB{ming.DB{config.RedisServer, config.RedisPassword}}
	records, err := db.GetAllRecords()
	if err != nil {
		return
	}
}

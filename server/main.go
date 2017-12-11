package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/northbright/pathhelper"
)

var (
	serverRoot    = ""
	templatesPath = ""
	staticPath    = ""
	configFile    = ""
	config        Config
)

// Config represents the app settings.
type Config struct {
	ServerAddr    string   `json:"server_addr"`
	RedisServer   string   `json:"redis_server"`
	RedisPassword string   `json:"redis_password"`
	AdminAccount  string   `json:"admin_account"`
	AdminPassword string   `json:"admin_password"`
	ClosedNotices []string `json:"closed_notices"`
}

func main() {
	var (
		err        error
		authorized *gin.RouterGroup
	)

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	serverRoot, _ = pathhelper.GetCurrentExecDir()
	templatesPath = path.Join(serverRoot, "templates")
	staticPath = path.Join(serverRoot, "static")

	if err = loadConfig("config.json", &config); err != nil {
		return
	}

	r := gin.Default()

	// Serve Static files.
	r.Static("/static/", staticPath)

	// Load Templates.
	r.LoadHTMLGlob(fmt.Sprintf("%v/*", templatesPath))

	// Pages
	r.GET("/", home)
	authorized = r.Group("/", gin.BasicAuth(gin.Accounts{
		config.AdminAccount: config.AdminPassword,
	}))
	authorized.GET("/admin", admin)

	r.Run(config.ServerAddr)
}

func loadConfig(file string, config *Config) error {
	var (
		err        error
		buf        []byte
		currentDir string
	)

	currentDir, _ = pathhelper.GetCurrentExecDir()
	file = path.Join(currentDir, file)

	// Load Conifg
	if buf, err = ioutil.ReadFile(file); err != nil {
		return err
	}

	return json.Unmarshal(buf, &config)
}

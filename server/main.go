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
	redisAddr     = ":6379"
	redisPassword = ""
	serverRoot    = ""
	templatesPath = ""
	staticPath    = ""
	configFile    = ""
)

type Config struct {
	AdminAccount  string
	AdminPassword string
}

func main() {
	var err error
	var authorized *gin.RouterGroup
	buf := []byte{}
	config := Config{}

	r := gin.Default()

	serverRoot, _ = pathhelper.GetCurrentExecDir()
	templatesPath = path.Join(serverRoot, "templates")
	staticPath = path.Join(serverRoot, "static")

	configFile = path.Join(serverRoot, "config.json")

	// Load Conifg
	if buf, err = ioutil.ReadFile(configFile); err != nil {
		log.Printf("Load config file error: %v\n", err)
		goto end
	}

	if err = json.Unmarshal(buf, &config); err != nil {
		log.Printf("Parse config err: %v\n", err)
		goto end
	}

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

	//r.POST("/zb", postZB)

	r.Run(":8080")
end:
	if err != nil {
		log.Printf("main() error: %v\n", err)
		return
	}
}

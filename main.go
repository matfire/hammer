package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Apps map[string]App
}

type App struct {
	Name     string
	Path     string
	Commands []string
	Secret   string
}

func main() {
	var configPath string
	var config Config

	flag.StringVar(&configPath, "config", "./config.toml", `Path to config file (defaults to current dir's config.toml)`)
	flag.Parse()
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(config)
	r := gin.Default()
	r.GET("/up", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/trigger/:project", func(ctx *gin.Context) {
		project := ctx.Param("project")
		event := ctx.GetHeader("x-github-event")
		if _, ok := config.Apps[project]; ok {
			//TODO do the thing with the secret
			switch event {
			case "ping":
				ctx.String(200, "pong")
				break
			case "release":
				break
			}
		}
		ctx.String(200, project)
	})
	r.Run()
}

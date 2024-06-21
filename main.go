package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
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
		if projectConfig, ok := config.Apps[project]; ok {
			//TODO do the thing with the secret
			signature := ctx.GetHeader("X-Hub-Signature-256")
			payload, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				//TODO do something here
			}
			mac := hmac.New(sha256.New, []byte(projectConfig.Secret))
			mac.Write(payload)
			if !hmac.Equal([]byte(signature[7:]), mac.Sum(nil)) {
				ctx.String(500, "invalid secret")
			}
			switch event {
			case "ping":
				break
			case "release":
				break
			}
		}
	})
	r.Run()
}

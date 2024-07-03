package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/matfire/hammer/git"
	"github.com/matfire/hammer/types"
)

func main() {
	var configPath string
	var config types.Config

	flag.StringVar(&configPath, "config", "./config.toml", `Path to config file (defaults to current dir's config.toml)`)
	flag.Parse()
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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
				ctx.String(500, "could not parse body")
				return
			}
			mac := hmac.New(sha256.New, []byte(projectConfig.Secret))
			mac.Write(payload)
			macValue := hex.EncodeToString(mac.Sum(nil))
			if !hmac.Equal([]byte(signature[7:]), []byte(macValue)) {
				ctx.String(500, "invalid secret")
				return
			}
			switch event {
			case "ping":
				break
			case "release":
				var releasePayload types.GithubReleasePayload
				if err := json.Unmarshal(payload, &releasePayload); err != nil {
					ctx.String(500, "cannot parse body data")
					return
				}
				fmt.Println(releasePayload)
				git.Pull(config, projectConfig, releasePayload)
				break
			}
			ctx.String(200, "ok")
		}
	})
	r.Run()
}

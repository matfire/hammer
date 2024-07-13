package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/matfire/hammer/exex"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/matfire/hammer/git"
	"github.com/matfire/hammer/types"
)

func main() {
	var configPath string
	var config types.Config
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	flag.StringVar(&configPath, "config", "./config.toml", `Path to config file (defaults to current dir's config.toml)`)
	flag.Parse()
	logger.Info("parsed flags")
	logger.Info("decoding toml file")
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		logger.Error("failed decoding toml file")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logger.Info("decoded toml file")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/up", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/trigger/:project", func(ctx *gin.Context) {
		project := ctx.Param("project")
		event := ctx.GetHeader("x-github-event")
		if projectConfig, ok := config.Apps[project]; ok {
			logger.Info("triggering event", "project", project, "event", event)
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
			case "release":
				var releasePayload types.GithubReleasePayload
				if err := json.Unmarshal(payload, &releasePayload); err != nil {
					ctx.String(500, "cannot parse body data")
					return
				}
				git.Pull(config, projectConfig, releasePayload)
				for i := 0; i < len(projectConfig.Commands); i++ {
					logger.Info("executing command", "project", project, "command", projectConfig.Commands[i], "index", i)
					err = exex.Exec(projectConfig.Commands[i], projectConfig.Path)
					if err != nil {
						logger.Error("failed executing command", "project", project, "command", projectConfig.Commands[i], "index", i)
						ctx.String(500, "failed to execute command number "+strconv.Itoa(i))
						break
					}
					logger.Info("finished executing command", "project", project, "command", projectConfig.Commands[i], "index", i)
				}
				break
			default:
				ctx.String(500, "unsupported event")
			}
		} else {
			logger.Error("failed to process event for project", "project", project)
		}
	})
	err = r.Run()
	if err != nil {
		panic("could not run server")
	}
}

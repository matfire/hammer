package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/matfire/hammer/server"
	"github.com/matfire/hammer/types"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	var configPath string
	var port int
	var debug bool
	var logFile string
	var config types.Config

	flag.StringVar(&configPath, "config", "./config.toml", `Path to config file (defaults to current dir's config.toml)`)
	flag.IntVar(&port, "port", 8080, "port to run the webserver on (defaults to 8080)")
	flag.BoolVar(&debug, "debug", false, "turns on debug mode for the web server (defaults to off)")
	flag.StringVar(&logFile, "log", "stdout", "where to write logs (defaults to stdout)")
	flag.Parse()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if logFile != "stdout" {
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		logger = slog.New(slog.NewJSONHandler(file, nil))
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	logger.Info("parsed flags")
	logger.Info("decoding toml file")
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		logger.Error("failed decoding toml file")
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logger.Info("decoded toml file")
	instance := server.NewServer(&config, logger)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), instance)
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	if err != nil {
		panic("could not start server")
	}
}

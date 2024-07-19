package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/matfire/hammer/server"
	"github.com/matfire/hammer/types"
	"log/slog"
	"net/http"
	"os"
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
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logger.Info("decoded toml file")
	instance := server.NewServer(&config, logger)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), instance)
	if err != nil {
		panic("could not start server")
	}
}

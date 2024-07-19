package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/matfire/hammer/exec"
	"github.com/matfire/hammer/git"
	"github.com/matfire/hammer/types"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func handleTrigger(w http.ResponseWriter, r *http.Request, config *types.Config, logger *slog.Logger) {
	project := r.PathValue("project")
	event := r.Header.Get("x-github-event")
	if projectConfig, ok := config.Apps[project]; ok {
		logger.Info("triggering event", "project", project, "event", event)
		signature := ctx.GetHeader("X-Hub-Signature-256")
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("could not parse body"))
			return
		}
		mac := hmac.New(sha256.New, []byte(projectConfig.Secret))
		mac.Write(payload)
		macValue := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(signature[7:]), []byte(macValue)) {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("could not validate secret"))
			return
		}
		switch event {
		case "release":
			var releasePayload types.GithubReleasePayload
			if err := json.Unmarshal(payload, &releasePayload); err != nil {
				w.WriteHeader(500)
				_, _ = w.Write([]byte("could not parse body data"))
				return
			}
			git.Pull(projectConfig, releasePayload)
			for i := 0; i < len(projectConfig.Commands); i++ {
				logger.Info("executing command", "project", project, "command", projectConfig.Commands[i], "index", i)
				err = exec.Exec(projectConfig.Commands[i], projectConfig.Path)
				if err != nil {
					logger.Error("failed executing command", "project", project, "command", projectConfig.Commands[i], "index", i)
					w.WriteHeader(500)
					_, _ = w.Write([]byte("failed to execute command number " + strconv.Itoa(i)))
					break
				}
				logger.Info("finished executing command", "project", project, "command", projectConfig.Commands[i], "index", i)
			}
			break
		default:
			w.WriteHeader(500)
			_, _ = w.Write([]byte("unsupported event"))
		}
	} else {
		logger.Error("failed to process event for project", "project", project)
	}
}

func NewServer(config *types.Config, logger *slog.Logger) *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/trigger/{project}", func(w http.ResponseWriter, r *http.Request) {
		handleTrigger(w, r, config, logger)
	})
	server.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("'status': 'OK'"))
	})
	return server
}

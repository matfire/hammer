package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/matfire/hammer/types"
)

func Pull(config types.Config, app types.App, releasePayload types.GithubReleasePayload) {
	auth := &http.BasicAuth{
		Username: config.Username,
		Password: config.Password,
	}
	repo, err := git.PlainOpen(app.Path)
	if err != nil {
		panic("path is not a valid repo")
	}
	fmt.Println(repo)
	workTree, err := repo.Worktree()
	if err != nil {
		panic("could not get worktree")
	}
	err = workTree.Pull(&git.PullOptions{Auth: auth})
	if err != nil {
		panic("could not pull")
	}
}

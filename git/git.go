package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5"
	"github.com/matfire/hammer/types"
)

func Pull(app types.App, releasePayload types.GithubReleasePayload) {
	repo, err := git.PlainOpen(app.Path)
	if err != nil {
		panic("path is not a valid repo")
	}
	workTree, err := repo.Worktree()
	if err != nil {
		panic("could not get worktree")
	}
	err = workTree.Pull(&git.PullOptions{})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		fmt.Println(err)
		panic("could not pull")
	}
	err = workTree.Checkout(&git.CheckoutOptions{Branch: plumbing.NewTagReferenceName(releasePayload.Release.TagName)})
	if err != nil {
		fmt.Println(err)
		panic("could not checkout specified tag")
	}
}

package sync

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"os"
)

type Git struct {
	repoURL string
}

func NewGit(repoURL string) *Git {
	return &Git{repoURL: repoURL}
}

func (g *Git) Clone(dst string) error {
	_, err := git.PlainClone(dst, false, &git.CloneOptions{
		URL:      g.repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("cloning repo: %v", err)
	}
	return nil
}

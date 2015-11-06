package skel

import (
	"errors"
	"os"

	"github.com/roboll/skel/skel/github"
)

type Source interface {
	DataLocation() (*string, error)
}

type DirectorySource struct {
	Path string
}

type GithubReleaseSource struct {
	Tag   string
	Owner string
	Repo  string
	Name  string
}

func (s *DirectorySource) DataLocation() (*string, error) {
	//trim trailing slash if it has one
	for len(s.Path) > 0 && s.Path[len(s.Path)-1] == '/' {
		s.Path = s.Path[0 : len(s.Path)-1]
	}

	info, err := os.Stat(s.Path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("DirectorySource: Path must be a directory.")
	}
	return &s.Path, nil
}

func (g *GithubReleaseSource) validate() error {
	if len(g.Tag) == 0 {
		return errors.New("GithubReleaseSource: Tag is required.")
	}
	if len(g.Owner) == 0 {
		return errors.New("GithubReleaseSource: Owner is required.")
	}
	if len(g.Repo) == 0 {
		return errors.New("GithubReleaseSource: Repo is required.")
	}
	if len(g.Name) == 0 {
		return errors.New("GithubReleaseSource: Name is required.")
	}
	return nil
}

func (g *GithubReleaseSource) DataLocation() (*string, error) {
	if err := g.validate(); err != nil {
		return nil, err
	}

	gh := github.Github{Token: os.Getenv("GITHUB_TOKEN")}
	dl, err := gh.DownloadRelease(g.Owner, g.Repo, g.Name, g.Tag)
	if err != nil {
		return nil, err
	}
	return &dl, nil
}

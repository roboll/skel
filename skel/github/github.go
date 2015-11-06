package github

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	gh "github.com/google/go-github/github"
	"github.com/pivotal-golang/archiver/extractor"

	"golang.org/x/oauth2"
)

const (
	archiveTemplate = "https://%sapi.github.com/repos/%s/%s/releases/assets/%d"
)

type Github struct {
	Token string
}

func (g *Github) DownloadRelease(owner, repo, name, tag string) error {
	filename, err := g.downloadRelease(owner, repo, name, tag)
	if err != nil {
		return err
	}
	workdir, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Rename(*filename, path.Join(workdir, name+".tar.gz"))
	return nil
}

func (g *Github) downloadRelease(owner, repo, name, tag string) (*string, error) {
	url, err := g.getDownloadUrl(owner, repo, name+".tar.gz", tag)
	if err != nil {
		return nil, err
	}
	filename := path.Join(os.TempDir(), name+".tar.gz")

	log.Printf("downloading asset to %s\n", filename)
	err = g.download(*url, filename)
	if err != nil {
		return nil, err
	}
	return &filename, nil
}

func (g *Github) DownloadAndExtractRelease(owner, repo, name, tag string) (*string, error) {
	filename, err := g.downloadRelease(owner, repo, name, tag)
	if err != nil {
		return nil, err
	}
	if filename == nil {
		return nil, errors.New("Filename came back nil. This is a bug.")
	}
	outdir := path.Join(os.TempDir(), name)

	log.Println("extracting asset")
	err = extractor.NewTgz().Extract(*filename, outdir)
	if err != nil {
		return nil, err
	}

	return &outdir, nil
}

func (g *Github) getDownloadUrl(owner, repo, name, tag string) (*string, error) {
	client := gh.NewClient(client(g.Token))

	var rel *gh.RepositoryRelease
	var err error
	switch tag {
	case "latest":
		rel, _, err = client.Repositories.GetLatestRelease(owner, repo)
	default:
		rel, _, err = client.Repositories.GetReleaseByTag(owner, repo, tag)
	}
	if err != nil {
		return nil, err
	}

	var asset *gh.ReleaseAsset
	for _, a := range rel.Assets {
		if *a.Name == name {
			asset = &a
			break
		}
	}

	if asset == nil {
		return nil, fmt.Errorf("no asset named %s", name)
	}

	var token string
	if len(g.Token) > 0 {
		token = g.Token + "@"
	} else {
		token = ""
	}
	url := fmt.Sprintf(archiveTemplate, token, owner, repo, *asset.ID)
	return &url, nil
}

func (g *Github) download(url, filename string) error {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/octet-stream")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	var content bytes.Buffer
	content.ReadFrom(resp.Body)

	err = ioutil.WriteFile(filename, content.Bytes(), os.ModePerm)
	return err
}

func client(token string) *http.Client {
	if len(token) == 0 {
		return http.DefaultClient
	}
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	return oauth2.NewClient(oauth2.NoContext, src)
}

package skel

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/codegangsta/cli"
	"github.com/roboll/skel/skel/github"
	"github.com/roboll/skel/skel/tmpl"
)

var release string

func Run(args []string) {
	app := cli.NewApp()

	app.Name = "skel"
	app.Usage = "https://github.com/roboll/skel"

	if len(release) == 0 {
		release = "HEAD"
	}
	app.Version = release

	app.Flags = flags
	app.Action = run

	app.Run(args)
}

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "data",
		Usage: "yaml file with template data",
		Value: "data.yaml",
	},
	cli.StringFlag{
		Name:  "dest",
		Usage: "dest dir",
		Value: "./",
	},
	cli.StringFlag{
		Name:  "src",
		Usage: "src dir - takes presidence over gh- options",
	},
	cli.StringFlag{
		Name:  "skel",
		Usage: "release artifact name",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "name",
		Value: "",
	},
	cli.StringFlag{
		Name:  "gh-tag",
		Usage: "release tag",
		Value: "latest",
	},
	cli.StringFlag{
		Name:   "gh-owner",
		Usage:  "github repo owner",
		Value:  "roboll",
		EnvVar: "SKEL_OWNER",
	},
	cli.StringFlag{
		Name:   "gh-repo",
		Usage:  "github repo",
		Value:  "skel",
		EnvVar: "SKEL_REPO",
	},
	cli.StringFlag{
		Name:   "gh-token",
		Usage:  "github api token",
		EnvVar: "GITHUB_TOKEN",
	},
}

func run(c *cli.Context) {
	var prefix string

	dest := c.String("dest")
	src := c.String("src")

	if len(src) == 0 {
		//github
		name := c.String("skel")
		if len(name) == 0 {
			log.Fatal("skel is required")
		}

		tag := c.String("gh-tag")
		owner := c.String("gh-owner")
		repo := c.String("gh-repo")
		token := c.String("gh-token")

		gh := github.Github{Token: token}
		dl, err := gh.DownloadRelease(owner, repo, name, tag)
		if err != nil {
			log.Fatal(err)
		}
		src = dl
		prefix = os.TempDir() + name
	} else {
		prefix = src
	}

	log.Printf("prefix is %s", prefix)

	var data map[string]string

	defaultData, err := ioutil.ReadFile(path.Join(src, "skel.yaml"))
	if err != nil {
		log.Println("skel: failed to load default data, ignoring.")
	} else {
		err = yaml.Unmarshal(defaultData, &data)
		if err != nil {
			log.Println("skel: failed to unmarshal skel.yaml, ignoring.")
		}
	}

	filepath := c.String("data")
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("skel: couldn't read data file %s: %s\n", filepath, err)
	}

	err = yaml.Unmarshal(content, &data)
	if err != nil {
		log.Fatalf("skel: unable to parse %s: %s\n", filepath, err)
	}

	name := c.String("name")
	data["name"] = name
	err = tmpl.Template(src, dest, name, prefix, []string{"skel.yaml"}, data)
	if err != nil {
		log.Fatalf("skel: error templating %s: %s\n", filepath, err)
	}
}

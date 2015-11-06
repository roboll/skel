package skel

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/codegangsta/cli"
)

func RunCli(args []string) {
	app := cli.NewApp()

	app.Name = "skel"
	app.Usage = "https://github.com/roboll/skel"

	if len(release) == 0 {
		release = "HEAD"
	}
	app.Version = release

	app.Flags = flags
	app.Action = clirun

	app.Run(args)
}

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "data",
		Usage: "yaml file with template data",
		Value: "data.yaml",
	},
	cli.BoolTFlag{
		Name:  "open-editor",
		Usage: "open editor with data before templating",
	},
	cli.StringFlag{
		Name:  "dest",
		Usage: "dest dir",
	},
	cli.StringFlag{
		Name:  "src",
		Usage: "src directory - takes presidence over gh- options",
	},
	cli.StringFlag{
		Name:  "skel",
		Usage: "release artifact name",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "name",
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
}

func clirun(c *cli.Context) {
	config := &Config{
		Dest: c.String("dest"),
		Name: c.String("name"),

		Data:       make(map[string]string),
		OpenEditor: c.BoolT("open-editor"),
	}

	src := c.String("src")
	if len(src) > 0 {
		config.Source = &DirectorySource{Path: src}
	} else {
		config.Source = &GithubReleaseSource{
			Tag:   c.String("gh-tag"),
			Owner: c.String("gh-owner"),
			Repo:  c.String("gh-repo"),
			Name:  c.String("skel"),
		}
	}

	data := c.String("data")
	if len(data) > 0 {
		info, err := os.Stat(data)
		if err != nil || info.IsDir() {
			log.Printf("Couldn't open %s. Skipping.\n", data)
		} else {
			content, err := ioutil.ReadFile(data)
			if err != nil {
				log.Printf("Failed to read data file. Skipping. Error was: %s\n", err)
			}
			err = yaml.Unmarshal(content, &config.Data)
			if err != nil {
				log.Printf("Failed to unmarshall data file: %s.\n", err)
			}
		}
	}

	if err := Run(config); err != nil {
		log.Fatal(err)
	}
}

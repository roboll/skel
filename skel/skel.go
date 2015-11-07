package skel

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/roboll/skel/skel/tmpl"

	"gopkg.in/yaml.v2"
)

var release string

type Config struct {
	Source Source
	Dest   string
	Name   string

	Data       map[string]string
	OpenEditor bool
}

func (c *Config) validate() error {
	if c.Source == nil {
		return errors.New("Source may not be nil.")
	}
	if len(c.Dest) == 0 {
		workdir, err := os.Getwd()
		if err != nil {
			return err
		}
		c.Dest = workdir
		log.Printf("skel: Dest was empty. Using %s.\n", c.Dest)
	}
	if len(c.Name) == 0 {
		idx := strings.LastIndex(c.Dest, "/")
		c.Name = c.Dest[idx+1 : len(c.Dest)]
		log.Printf("skel: Name was empty. Using %s.\n", c.Name)
	}
	if c.Data == nil {
		c.Data = make(map[string]string)
	}
	return nil
}

func Run(config *Config) error {
	if err := config.validate(); err != nil {
		return err
	}

	src, err := config.Source.GetLocation()
	if err != nil {
		return err
	}
	if src == nil {
		return errors.New("skel: Source returned nil pointer, this is a bug.")
	}
	defer config.Source.Cleanup()

	prefix := *src

	var data map[string]string = make(map[string]string)

	defaultData, err := ioutil.ReadFile(path.Join(*src, "skel.yaml"))
	if err != nil {
		log.Println("skel: Failed to load default data from template's skel.yaml, ignoring.")
	} else {
		err = yaml.Unmarshal(defaultData, &data)
		if err != nil {
			log.Println("skel: Failed to unmarshal template's skel.yaml, ignoring.")
		}
	}

	for key, val := range config.Data {
		data[key] = val
	}
	data["name"] = config.Name

	if config.OpenEditor {
		err := doEditData(config, data)
		if err != nil {
			return err
		}
	} else {
		log.Println("skel: OpenEditor was false, not editing.")
	}

	if err := tmpl.Template(*src, config.Dest, data["name"], prefix, []string{"skel.yaml"}, data); err != nil {
		return err
	}

	return nil
}

func doEditData(config *Config, data map[string]string) error {
	path := path.Join(os.TempDir(), "skel-data.yaml")
	tmp, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Failed to open temp file for editing: %s.", err)
	}
	out, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("Failed to write data to tempfile: %s.", err)
	}
	if _, err := tmp.Write(out); err != nil {
		return fmt.Errorf("Failed to write data to tempfile: %s.", err)
	}
	tmp.Close()

	cmd := exec.Command(os.Getenv("EDITOR"), path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to open $EDITOR: %s.", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Failed to wait for $EDITOR to close: %s.", err)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read data file after editing: %s.", err)
	}
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		return fmt.Errorf("Failed to unmarshall yaml content after editing: %s.", err)
	}
	return nil
}

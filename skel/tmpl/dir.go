package tmpl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

func Template(src, dest, name, prefix string, ignore []string, data interface{}) error {
	d := &holder{data: data, dest: dest, name: name, prefix: prefix, ignore: ignore}
	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return fmt.Errorf("template: failed to create %s: %s\n", dest, err)
	}
	return filepath.Walk(src, d.handle)
}

type holder struct {
	data   interface{}
	dest   string
	name   string
	prefix string
	ignore []string
}

func (h *holder) handle(loc string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("template: error: %s\n", err)
		return nil
	}

	for _, str := range h.ignore {
		if info.Name() == str {
			log.Printf("template: found ignore match %s, skipping.\n", str)
			return nil
		}
	}

	dest := path.Join(h.dest, strings.Replace(loc, h.prefix, "", -1))
	dest = strings.Replace(dest, "skel", h.name, -1)

	if info.IsDir() {
		// dir: create empty dir
		log.Printf("template: creating dir %s\n", dest)

		info, err := os.Stat(dest)
		if err == nil && info.IsDir() {
			log.Printf("template: %s already exists, skipping.", dest)
		} else {
			err = os.Mkdir(dest, os.ModePerm)
			if err != nil {
				log.Printf("template: error creating %s: %s\n", dest, err)
			}
		}
	} else {
		// file: template file to dest
		log.Printf("template: processing %s", loc)
		tmpl, err := template.New(filepath.Base(loc)).Delims("{{{", "}}}").ParseFiles(loc)
		if err != nil {
			log.Printf("template: error loading template %s: %s\n", loc, err)
		}
		if tmpl == nil {
			log.Printf("template: error while processing template %s.\n", loc)
		}

		file, err := os.Create(dest)
		if err != nil {
			log.Printf("template: error creating %s: %s\n", path.Join(h.dest, loc), err)
		}
		defer file.Close()

		if err := file.Chmod(info.Mode()); err != nil {
			log.Printf("template: failed to set permissions on %s: %s", file.Name(), err)
		}

		writer := bufio.NewWriter(file)
		err = tmpl.Execute(writer, h.data)
		defer writer.Flush()
		if err != nil {
			log.Printf("template: error writing template %s: %s\n", path.Join(h.dest, loc), err)
		}
	}
	return nil
}

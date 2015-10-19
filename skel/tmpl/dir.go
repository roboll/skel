package tmpl

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Template(src, dest, prefix string, ignore []string, data interface{}) error {
	d := &holder{data: data, dest: dest, prefix: prefix, ignore: ignore}
	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return fmt.Errorf("template: failed to create %s: %s\n", dest, err)
	}
	return filepath.Walk(src, d.handle)
}

type holder struct {
	data   interface{}
	dest   string
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
			log.Printf("found ignore match %s, skipping.\n", str)
			return nil
		}
	}

	if info.IsDir() {
		// dir: create empty dir
		dest := path.Join(h.dest, strings.Replace(loc, h.prefix, "", -1))
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
		tmpl, err := template.ParseFiles(loc)
		if err != nil {
			log.Printf("template: error loading template %s: %s\n", loc, err)
		}

		dest := path.Join(h.dest, strings.Replace(loc, h.prefix, "", -1))
		file, err := os.Create(dest)
		if err != nil {
			log.Printf("template: error creating %s: %s\n", path.Join(h.dest, loc), err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		err = tmpl.Execute(writer, h.data)
		defer writer.Flush()
		if err != nil {
			log.Printf("template: error writing template %s: %s\n", path.Join(h.dest, loc), err)
		}
	}
	return nil
}

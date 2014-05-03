package jig

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type JigSpec struct {
	Pre     []string
	Build   []string
	Post    []string
	Image   string
	Name    string
	Workdir string
	Environ map[string]string
	Mount   string
	Jigfile *Jigfile `json:"-"`
}

const (
	JIGDIR string = ".jig"
)

func asByteLines(s []string) []byte {
	return []byte(strings.Join(s, "\n"))
}

func (spec *JigSpec) ConfigDir() (string, error) {
	return ensureDirectory(spec.Jigfile.Path, JIGDIR, spec.Name)
}

func (spec *JigSpec) updateIfDiff(data []byte, filename string) error {
	cfgdir, err := spec.ConfigDir()
	if err != nil {
		return err
	}
	fullname := filepath.Join(cfgdir, filename)
	existing, err := ioutil.ReadFile(fullname)
	if err != nil || md5.Sum(data) != md5.Sum(existing) {
		if err := ioutil.WriteFile(fullname, data, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (spec *JigSpec) CreateScripts() error {
	mkdir := fmt.Sprintf("mkdir -p %s && chown -R jig:jig /mnt/jig", spec.Mount)
	spec.Pre = append([]string{mkdir}, spec.Pre...)
	if err := spec.updateIfDiff(asByteLines(spec.Pre), "pre"); err != nil {
		return err
	}
	if err := spec.updateIfDiff(asByteLines(spec.Build), "build"); err != nil {
		return err
	}
	if err := spec.updateIfDiff(asByteLines(spec.Post), "post"); err != nil {
		return err
	}
	return nil
}
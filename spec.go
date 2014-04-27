package jig

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

type JigSpec struct {
	Pre     []string
	Build   []string
	Post    []string
	Output  []string
	Image   string
	Name    string
	Jigfile *Jigfile
}

const (
	JIGDIR string = ".jig"
)

func (spec *JigSpec) ConfigDir() (string, error) {
	return ensureDirectory(spec.Jigfile.Path, JIGDIR, spec.Name)
}

func (spec *JigSpec) CreateScripts() error {
	cfgdir, err := spec.ConfigDir()
	if err != nil {
		return err
	}
	pre := filepath.Join(cfgdir, "pre")
	predata := strings.Join(spec.Pre, "\n")
	if err := ioutil.WriteFile(pre, []byte(predata), 0755); err != nil {
		return err
	}
	build := filepath.Join(cfgdir, "build")
	builddata := strings.Join(spec.Build, "\n")
	if err := ioutil.WriteFile(build, []byte(builddata), 0755); err != nil {
		return err
	}
	post := filepath.Join(cfgdir, "post")
	postdata := strings.Join(spec.Post, "\n")
	if err := ioutil.WriteFile(post, []byte(postdata), 0755); err != nil {
		return err
	}
	return nil
}
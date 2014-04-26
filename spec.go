package jig

import (
	"os"
	"path/filepath"
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
	p := filepath.Join(spec.Jigfile.Path, JIGDIR, spec.Name)
	if err := os.MkdirAll(p, 0664); err != nil {
		return "", err
	}
	return p, nil
}
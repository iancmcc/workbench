package jig

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const (
	JIGFILE = "Jigfile"
)

type Jigfile struct {
	Path  string
	Specs map[string]*JigSpec
}

func ParseJigfile(r io.Reader) (*Jigfile, error) {
	jf := &Jigfile{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(&jf.Specs); err != nil && err != io.EOF {
		return nil, err
	}
	for name, spec := range jf.Specs {
		spec.Jigfile = jf
		spec.Name = name
	}
	return jf, nil
}

func ParseJigfilePath(p string) (*Jigfile, error) {
	var (
		dir   *os.File
		f     io.ReadCloser
		fname string = p
		err   error
	)
	if dir, err = os.Open(p); err != nil {
		return nil, err
	}
	defer dir.Close()
	stat, err := dir.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		fname = filepath.Join(p, JIGFILE)
		if f, err = os.Open(fname); err != nil {
			return nil, err
		}
		defer f.Close()
	} else {
		f = dir
	}
	jf, err := ParseJigfile(f)
	if err != nil {
		return nil, err
	}
	jf.Path = dir.Name()
	return jf, nil
}

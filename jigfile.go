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
	Builds map[string]Build
}

type Build struct {
	Pre    []string
	Build  []string
	Post   []string
	Output []string
	Image  string
}

func ParseJigfile(r io.Reader) (*Jigfile, error) {
	jf := &Jigfile{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(&jf.Builds); err != nil && err != io.EOF {
		return nil, err
	}
	return jf, nil
}

func ParseJigfilePath(p string) (*Jigfile, error) {
	var (
		dir *os.File
		f   io.ReadCloser
		err error
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
		if f, err = os.Open(filepath.Join(p, JIGFILE)); err != nil {
			return nil, err
		}
		defer f.Close()
	} else {
		f = dir
	}
	return ParseJigfile(f)
}

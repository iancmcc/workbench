package jig

import (
	"encoding/json"
	"io"
	"os"
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

func ParseJigfileFromFile(p string) (*Jigfile, error) {
	var (
		f   io.ReadCloser
		err error
	)
	if f, err = os.Open(p); err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseJigfile(f)
}

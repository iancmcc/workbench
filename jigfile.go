package jig

import (
	"encoding/json"
	"io"
)

type Jigfile struct {
	Builds map[string]Build
}

type Build struct {
	Pre    []string
	Build  []string
	Post   []string
	Output []string
}

func ParseJigfile(r io.Reader) (*Jigfile, error) {
	jf := &Jigfile{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(&jf.Builds); err != nil && err != io.EOF {
		return nil, err
	}
	return jf, nil
}

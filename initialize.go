package jig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Initialize(dir, name, image string, pre, build, post []string) error {
	var (
		jigfile string
		err     error
	)
	if stat, err := os.Stat(dir); err != nil || !stat.Mode().IsDir() {
		return fmt.Errorf("%s isn't a valid directory.")
	}
	jigfile = filepath.Join(dir, JIGFILE)
	if _, err := os.Stat(jigfile); err == nil {
		return fmt.Errorf("Jigfile already exists.")
	}
	specs := map[string]*JigSpec{}
	specs[name] = &JigSpec{
		Image:   image,
		Pre:     pre,
		Build:   build,
		Post:    post,
		Environ: make(map[string]string),
		Mount:   WORKBENCH,
		Workdir: WORKBENCH,
	}
	data, err := json.MarshalIndent(specs, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile(jigfile, data, 0664)
	return nil
}

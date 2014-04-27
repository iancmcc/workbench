package jig

import (
	"os"
	"path/filepath"
)

func ensureDirectory(path ...string) (string, error) {
	p := filepath.Join(path...)
	if err := os.MkdirAll(p, 0775); err != nil {
		return "", err
	}
	return p, nil
}

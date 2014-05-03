package jig

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
)

func ensureDirectory(path ...string) (string, error) {
	p := filepath.Join(path...)
	if err := os.MkdirAll(p, 0775); err != nil {
		return "", err
	}
	return p, nil
}

func randstring(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	en := base64.StdEncoding
	d := make([]byte, en.EncodedLen(len(b)))
	en.Encode(d, b)
	return strings.ToLower(string(d))
}

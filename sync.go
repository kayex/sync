package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Target interface {
	Write(r io.Reader, key string) error
}

type Sync struct {
	l *log.Logger
	t Target
}

func NewSync(l *log.Logger, t Target) *Sync {
	return &Sync{l:l, t: t}
}

func (s *Sync) Sync(dst, src string) error {
	files, err := buildFileIndex(src)
	if err != nil {
		return fmt.Errorf("failed listing files: %v", err)
	}

	for _, f := range files {
		s.l.Printf("syncing %v", f)

		in, err := os.Open(f)
		if err != nil {
			return fmt.Errorf("failed opening source file: %v", err)
		}

		if !strings.HasSuffix(dst, "/") {
			dst = dst + "/"
		}
		path := dst + f
		err = s.t.Write(in, path)
		if err != nil {
			return fmt.Errorf("failed copying %v: %v", f, err)
		}
	}

	return nil
}

func buildFileIndex(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

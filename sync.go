package sync

import (
	"fmt"
	"github.com/kr/fs"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Target interface {
	Write(r io.Reader, path string) error
	Exists(dir string) (bool, error)
	Remove(path string) error
	RemoveDir(dir string) error
	CreateDirAll(path string) error
	Walk(root string) *fs.Walker
}

type Sync struct {
	l *log.Logger
	t Target
}

func NewSync(l *log.Logger, t Target) *Sync {
	return &Sync{l:l, t: t}
}

var sep = string(os.PathSeparator)

func (s *Sync) Sync(dst, src string) error {
	s.l.Printf("building file index...")
	start := time.Now()
	files, err := buildFileIndex(src)
	if err != nil {
		return fmt.Errorf("failed listing files: %v", err)
	}
	indexTime := time.Since(start)

	s.l.Printf("syncing...")

	start = time.Now()
	err = s.clear(dst)
	if err != nil {
		return fmt.Errorf("failed deleting files on target: %v", err)
	}
	clearTime := time.Since(start)

	if !strings.HasSuffix(dst, sep) {
		dst = dst + sep
	}
	if !strings.HasSuffix(src, sep) {
		src = src + sep
	}

	start = time.Now()
	for _, f := range files {
		path := filepath.Join(dst, strings.TrimPrefix(f, src))
		s.l.Printf("CREATE %v", path)

		in, err := os.Open(f)
		if err != nil {
			return fmt.Errorf("failed opening source file: %v", err)
		}

		err = s.createDir(path)
		if err != nil {
			return fmt.Errorf("failed creating directories: %v", err)
		}

		err = s.t.Write(in, path)
		if err != nil {
			return fmt.Errorf("failed copying %v: %v", f, err)
		}
	}
	addTime := time.Since(start)

	s.l.Println()
	s.l.Printf("Successfully synced %d files.", len(files))
	s.l.Printf("INDEX:  %s", indexTime)
	s.l.Printf("DELETE: %s", clearTime)
	s.l.Printf("CREATE: %s", addTime)

	return nil
}

func (s *Sync) clear(dir string) error {
	if !strings.HasSuffix(dir, sep) {
		dir = dir + sep
	}

	var directories []string

	walker := s.t.Walk(dir)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			return err
		}
		if walker.Path() == dir {
			continue
		}

		if walker.Stat().IsDir() {
			directories = append(directories, walker.Path())
			continue
		}

		s.l.Printf("DELETE %v", walker.Path())
		err := s.t.Remove(walker.Path())
		if err != nil {
			return fmt.Errorf("failed deleting %v: %v", walker.Path(), err)
		}
	}

	for i := len(directories) - 1; i >= 0; i-- {
		d := directories[i]
		err := s.t.RemoveDir(d)
		if err != nil {
			return fmt.Errorf("failed deleting %v: %v", d, err)
		}

	}

	return nil
}

func (s *Sync) createDir(path string) error {
	dir := filepath.Dir(path)
	exists, err := s.t.Exists(dir)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	err = s.t.CreateDirAll(dir)
	if err != nil {
		return err
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

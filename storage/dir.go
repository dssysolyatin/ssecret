package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var _ Storage = (*dirStorage)(nil)

type dirStorage struct {
	dir string
}

func NewDirStorage(dir string) (*dirStorage, error) {
	if dir[len(dir)-1] != os.PathSeparator {
		dir += string(os.PathSeparator)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &dirStorage{dir: dir}, nil
}

func (s *dirStorage) Delete(ctx context.Context, path string) error {
	return os.Remove(s.absPath(path))
}

func (s *dirStorage) Upload(ctx context.Context, path string, reader io.ReadSeeker) error {
	p := s.absPath(path)
	_, err := os.Stat(p)
	switch {
	case err == nil:
		return &ErrObjectAlreadyExists{Key: path}
	case err != nil && !os.IsNotExist(err):
		return err
	}

	if err = os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	reader.Seek(0, io.SeekStart)
	if _, err = io.Copy(file, reader); err != nil {
		return err
	}

	return nil
}

func (s *dirStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	return os.Open(s.absPath(path))
}

func (s *dirStorage) absPath(path string) string {
	return fmt.Sprintf("%s%s", s.dir, path)
}

func (s *dirStorage) Close() error {
	return nil
}

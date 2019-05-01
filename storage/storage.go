package storage

import (
	"context"
	"fmt"
	"io"
)

var _ error = (*ErrObjectAlreadyExists)(nil)

type Storage interface {
	Upload(ctx context.Context, path string, reader io.ReadSeeker) error
	Download(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
	Close() error
}

type ErrObjectAlreadyExists struct {
	Key string
}

func (e *ErrObjectAlreadyExists) Error() string {
	return fmt.Sprintf(`object "%s" is already exist`, e.Key)
}

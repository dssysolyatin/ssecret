package crypt

import (
	"context"
	"io"
)

type Cipher interface {
	Encrypt(ctx context.Context, reader io.Reader) (io.Reader, error)
	Decrypt(ctx context.Context, reader io.Reader) (io.Reader, error)
}

package crypt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

const (
	encryptBufferSize = 1024 * 1024 // 1MB
	decryptBufferSize = 1048604
)

var (
	_ Cipher    = (*aesCipher)(nil)
	_ io.Reader = (*lambdaReader)(nil)
)

type lambdaReaderFunc func(ctx context.Context, data []byte) ([]byte, error)
type lambdaReader struct {
	buf          []byte
	processedBuf []byte
	plainReader  io.Reader
	ctx          context.Context
	eof          bool
	lambda       lambdaReaderFunc
}

func (r *lambdaReader) Read(p []byte) (int, error) {
	writtenBytes := 0
	for {
		select {
		case <-r.ctx.Done():
			return writtenBytes, r.ctx.Err()
		default:
		}

		if len(r.processedBuf) == 0 {
			if r.eof {
				return writtenBytes, io.EOF
			}

			err := r.fillProcessedBuf()
			if err != nil {
				return writtenBytes, err
			}
		}

		if len(r.processedBuf) > len(p) {
			copy(p, r.processedBuf[0:len(p)])
			r.processedBuf = r.processedBuf[len(p):]

			writtenBytes += len(p)
			return writtenBytes, nil
		} else {
			copy(p, r.processedBuf)
			p = p[len(r.processedBuf):]
			writtenBytes += len(r.processedBuf)

			r.processedBuf = nil
		}
	}
}

func (r *lambdaReader) fillProcessedBuf() error {
	readBytes, err := r.plainReader.Read(r.buf)
	switch {
	case err == io.EOF:
		r.eof = true
	case err != nil:
		return err
	}

	if readBytes == 0 {
		return nil
	}

	r.processedBuf, err = r.lambda(r.ctx, r.buf[0:readBytes])
	return err
}

func newLambdaReader(ctx context.Context, plainReader io.Reader, bufferSize int, lambda lambdaReaderFunc) *lambdaReader {
	return &lambdaReader{buf: make([]byte, bufferSize), plainReader: plainReader, ctx: ctx, lambda: lambda}
}

// TODO: protect variable `aesKey`
type aesCipher struct {
	c cipher.AEAD
}

func NewAesCipher(aesKey []byte) (*aesCipher, error) {
	cip, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	c, err := cipher.NewGCM(cip)
	if err != nil {
		return nil, err
	}

	return &aesCipher{c: c}, nil
}

func (c *aesCipher) Encrypt(ctx context.Context, reader io.Reader) (io.Reader, error) {
	return newLambdaReader(ctx, reader, encryptBufferSize, func(ctx context.Context, data []byte) ([]byte, error) {
		nonce := make([]byte, c.c.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}

		seal := c.c.Seal(nonce, nonce, data, nil)
		return seal, nil
	}), nil
}

func (c *aesCipher) Decrypt(ctx context.Context, reader io.Reader) (io.Reader, error) {
	return newLambdaReader(ctx, reader, decryptBufferSize, func(ctx context.Context, data []byte) ([]byte, error) {
		nonceSize := c.c.NonceSize()

		if len(data) < nonceSize {
			return nil, fmt.Errorf("invalid data size: data size should be greater than %d (NonceSize)", nonceSize)
		}

		nonce, encryptedData := data[:nonceSize], data[nonceSize:]
		plainData, err := c.c.Open(nil, nonce, encryptedData, nil)
		if err != nil {
			return nil, err
		}

		return plainData, nil
	}), nil
}

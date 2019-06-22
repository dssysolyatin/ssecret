package crypt

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLambdaReader(t *testing.T) {
	var expectedContent []byte
	lambdaReader := newLambdaReader(context.Background(), io.LimitReader(rand.Reader, 8*1024*1024), 1024*1024, func(ctx context.Context, data []byte) ([]byte, error) {
		expectedContent = append(expectedContent, data...)
		return data, nil
	})

	var actualContent []byte
	buf := make([]byte, 1024*1024+512)
	for {
		n, err := lambdaReader.Read(buf)
		actualContent = append(actualContent, buf[0:n]...)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}

	require.Equal(t, len(expectedContent), len(actualContent))
	require.Equal(t, expectedContent, actualContent)
}

func TestEncrypt(t *testing.T) {
	aesKey, err := ioutil.ReadAll(io.LimitReader(rand.Reader, 32))
	require.NoError(t, err)

	c, err := NewAesCipher([]byte(aesKey))
	require.NoError(t, err)

	expectedData, err := ioutil.ReadAll(io.LimitReader(rand.Reader, 10*1024*1024))
	require.NoError(t, err)

	encryptedReader, err := c.Encrypt(context.Background(), bytes.NewReader(expectedData))
	require.NoError(t, err)

	decryptedReader, err := c.Decrypt(context.Background(), encryptedReader)
	require.NoError(t, err)

	actualData, err := ioutil.ReadAll(decryptedReader)
	require.NoError(t, err)

	require.Equal(t, expectedData, actualData)

}

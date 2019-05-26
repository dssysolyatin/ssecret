package integration

import (
	"bytes"
	"context"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/dssysolyatin/ssecret/storage"
	"github.com/stretchr/testify/require"
)

func testStorage(storage storage.Storage, t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	b := make([]byte, 1024)
	_, err := rand.Read(b)
	require.NoError(t, err)

	path := "test"
	err = storage.Upload(ctx, path, bytes.NewReader(b))
	require.NoError(t, err)

	reader, err := storage.Download(ctx, path)
	require.NoError(t, err)
	defer reader.Close()

	bd, err := ioutil.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, bd, b)

	require.NoError(t, storage.Delete(ctx, path))
}

func TestDirStorage(t *testing.T) {
	dir := os.Getenv("DIR_STORAGE_DIR_PATH")
	st, err := storage.NewDirStorage(dir)
	require.NoError(t, err)
	defer os.Remove(dir)

	testStorage(st, t)
}

func TestS3Storage(t *testing.T) {
	s3St, err := storage.NewS3Storage(&aws.Config{
		Region: aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("S3_STORAGE_ACCESS_KEY"),
			os.Getenv("S3_STORAGE_PRIVATE_KEY"),
			"",
		),
	}, os.Getenv("S3_STORAGE_BUCKET"))
	require.NoError(t, err)

	testStorage(s3St, t)
}

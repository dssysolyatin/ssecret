package storage

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
)

var (
	_ Storage = (*s3Storage)(nil)
)

type s3Storage struct {
	bucket string
	client *s3.S3
}

func (s *s3Storage) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Key:    aws.String(path),
		Bucket: aws.String(s.bucket),
	})

	return err
}

func NewS3Storage(config *aws.Config, bucket string) (*s3Storage, error) {
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	return &s3Storage{client: s3.New(sess), bucket: bucket}, nil

}

func (s *s3Storage) Upload(ctx context.Context, path string, reader io.ReadSeeker) error {
	_, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	switch err.(type) {
	case nil:
		return &ErrObjectAlreadyExists{Key: path}
	case awserr.Error:
		if err.(awserr.Error).Code() != "NotFound" {
			return err
		}
	default:
		return err
	}

	_, err = s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(s.bucket),
		Key:                  aws.String(path),
		Body:                 reader,
		StorageClass:         aws.String(s3.ObjectStorageClassIntelligentTiering),
		ACL:                  aws.String(s3.BucketCannedACLPrivate),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String(s3.ServerSideEncryptionAes256),
	})

	return err
}

func (s *s3Storage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	out, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (s *s3Storage) Close() error {
	return nil
}

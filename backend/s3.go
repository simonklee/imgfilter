package backend

import (
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

// S3 is an Amazon S3 implementation of the ImageBackend
type S3 struct {
	b *s3.Bucket
}

func NewS3(access, secret, region, bucket string) *S3 {
	auth := aws.Auth{
		AccessKey: access,
		SecretKey: secret,
	}

	s := s3.New(auth, aws.Regions[region])

	return &S3{
		b: s.Bucket(bucket),
	}
}

func (s *S3) ReadFile(filename string) ([]byte, error) {
	return s.b.Get(filename)
}

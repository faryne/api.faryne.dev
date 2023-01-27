package s3

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	serverConfig "github.com/faryne/api-server/config"
	"github.com/gofiber/fiber/v2"
	"io"
	"time"
)

type instance struct {
	Client *s3.Client
	Bucket string
	Region string
}

func New(bucket, region string) fiber.Storage {
	var i = instance{}
	// initialize s3 client
	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(serverConfig.Config.AWS.Key, serverConfig.Config.AWS.Secret, "")),
		config.WithRegion(region))
	client := s3.NewFromConfig(cfg)
	i.Client = client
	i.Bucket = bucket
	return &i
}

func (i *instance) Get(key string) ([]byte, error) {
	output, err := i.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &i.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	var b = new(bytes.Buffer)
	if _, err := io.Copy(b, output.Body); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (i *instance) Set(key string, val []byte, _ time.Duration) error {
	var reader = new(bytes.Buffer)
	reader.Write(val)
	if _, err := i.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &i.Bucket,
		Key:    &key,
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   reader,
	}); err != nil {
		return err
	}
	return nil
}

func (i *instance) Delete(key string) error {
	if _, err := i.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &i.Bucket,
		Key:    &key,
	}); err != nil {
		return err
	}
	return nil
}

func (_ *instance) Reset() error {
	return nil
}

func (_ *instance) Close() error {
	return nil
}

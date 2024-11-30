package config

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	Client *s3.Client
	Bucket string
}

func NewS3Config(accessKey, secretKey, region, bucket string) (*S3Config, error) {
	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion(region))

	if err != nil {
		return nil, fmt.Errorf("AWS Config 실패: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Config{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (s *S3Config) MakePresignURL(key string) (string, error) {
	presignClient := s3.NewPresignClient(s.Client)

	presignReq, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	}, s3.WithPresignExpires(time.Minute*15))

	if err != nil {
		return "", fmt.Errorf("URL 생성 실패: %v", err)
	}

	return presignReq.URL, nil
}

package adapters

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type StorageS3 struct {
	accessKey  string
	secretKey  string
	bucketName string
	endpoint   string
	cli        *s3.Client
}

func NewStorageS3(ctx context.Context, accessKey, secretKey, bucketName, endpoint string) (*StorageS3, error) {
	result := &StorageS3{
		accessKey:  accessKey,
		secretKey:  secretKey,
		bucketName: bucketName,
		endpoint:   endpoint,
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(result.customCredentialProvider()),
		config.WithEndpointResolver(aws.EndpointResolverFunc(result.customEndpointResolver(endpoint))),
		config.WithHTTPClient(customHTTPClient()),
	)
	if err != nil {
		return nil, err
	}

	result.cli = s3.NewFromConfig(cfg)

	return result, nil
}

func (s *StorageS3) Upload(ctx context.Context, folderName, fileName string, fileSize int64, fileReader io.Reader) (string, error) {
	key := fmt.Sprintf("%s/%s-%s", folderName, uuid.NewString(), fileName)
	_, err := s.cli.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Body:          fileReader,
		ContentLength: fileSize,
		Key:           aws.String(key),
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucketName, key), nil
}

func (s *StorageS3) Remove(ctx context.Context, filePath string) error {
	_, err := s.cli.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageS3) customCredentialProvider() aws.CredentialsProvider {
	return aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     s.accessKey,
			SecretAccessKey: s.secretKey,
			Source:          "customProvider",
		}, nil
	})
}

func (s *StorageS3) customEndpointResolver(url string) func(service, region string) (aws.Endpoint, error) {
	return func(service, region string) (aws.Endpoint, error) {
		if service == "S3" {
			return aws.Endpoint{
				URL:               url,
				HostnameImmutable: true,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown service: %s", service)
	}
}

func customHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func (s *StorageS3) BucketExists(bucketName string) (bool, error) {
	_, err := s.cli.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {

		log.Printf("Nope %v.", bucketName)
	} else {
		log.Printf("Bucket %v exists and you already own it.", bucketName)
	}

	return exists, err
}
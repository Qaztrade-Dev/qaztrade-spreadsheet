package adapters

import (
	"archive/zip"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
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

func (s *StorageS3) DownloadArchive(ctx context.Context, folderName string) (io.ReadCloser, domain.RemoveFunction, error) {
	tempDir, err := os.MkdirTemp("", "archive")
	if err != nil {
		return nil, nil, err
	}
	defer os.RemoveAll(tempDir)

	paginator := s3.NewListObjectsV2Paginator(s.cli, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String(folderName),
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, nil, err
		}

		for _, object := range output.Contents {
			if err := downloadFile(ctx, s.cli, s.bucketName, *object.Key, tempDir); err != nil {
				return nil, nil, err
			}
		}
	}

	zipfileName, err := archiveDir(filepath.Join(tempDir, folderName))
	if err != nil {
		return nil, nil, err
	}

	zipfileReader, err := os.Open(zipfileName)
	if err != nil {
		return nil, nil, err
	}

	remover := func() error {
		return os.Remove(zipfileName)
	}

	return zipfileReader, remover, nil
}

func downloadFile(ctx context.Context, client *s3.Client, bucket, key, tempDir string) error {
	// The file path needs to be a string
	// that includes the name of the file you're downloading.
	// Assuming you want to keep the same name as in S3, you might use:
	var (
		filePath = filepath.Join(tempDir, key)
		dirPath  = filepath.Dir(filePath)
	)

	// Make sure the directory exists
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory, %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file, %w", err)
	}

	defer file.Close()

	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	resp, err := client.GetObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to download file, %w", err)
	}

	defer resp.Body.Close()

	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to readFrom file, %w", err)
	}

	return nil
}

func archiveDir(dirPath string) (string, error) {
	zipfile, err := os.CreateTemp("", "archive-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create archive file: %w", err)
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip if it is a directory
		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		dirInfo, err := d.Info()
		if err != nil {
			return fmt.Errorf("failed to get info: %w", err)
		}

		header, err := zip.FileInfoHeader(dirInfo)
		if err != nil {
			return fmt.Errorf("failed to create zip header: %w", err)
		}

		// Using slash ensures that the file path is correctly formed for zip files
		header.Name = filepath.ToSlash(filepath.Join(filepath.Base(dirPath), path[len(dirPath)+1:]))

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip writer: %w", err)
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("failed to write file to archive: %w", err)
		}

		return nil
	})

	return zipfile.Name(), nil
}

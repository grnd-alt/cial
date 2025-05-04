package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileService struct {
	minioClient *minio.Client
	bucketName  string
}

func InitFileService(url string, accessKey string, secretKey string, bucketName string, env string) (*FileService, error) {
	minioClient, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: env == "production",
	})
	if err != nil {
		return nil, err
	}

	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}
	return &FileService{
		minioClient: minioClient,
		bucketName:  bucketName,
	}, nil
}

func (f *FileService) UploadFile(objectName string, file io.Reader, fileSize int64, mimeType string) (string, error) {
	_, err := f.minioClient.PutObject(
		context.Background(),
		f.bucketName,
		objectName,
		file,
		fileSize,
		minio.PutObjectOptions{ContentType: mimeType, CacheControl: "private, max-age=31536000"},
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", f.minioClient.EndpointURL(), f.bucketName, objectName), nil
}

func (f *FileService) GetFileUrl(objectName string) (string, error) {
	url, err := f.minioClient.PresignedGetObject(context.Background(), f.bucketName, objectName, time.Duration(1000)*time.Second, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

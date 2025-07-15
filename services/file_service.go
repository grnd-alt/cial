package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/cors"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileService struct {
	minioClient *minio.Client
	bucketName  string
	cache       FileUrlCache
}

type FileUrlCache interface {
	Get(objectName string) (string, error)
	Set(objectName string, url string, ttl time.Duration) error
}

func InitFileService(url string, accessKey string, secretKey string, bucketName string, fileUrlCache FileUrlCache, env string) (*FileService, error) {
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
	test := minioClient.SetBucketCors(context.Background(), bucketName, cors.NewConfig([]cors.Rule{
		{
			MaxAgeSeconds: 3600,
			AllowedHeader: []string{"*"},
			AllowedMethod: []string{"GET"}, // Allow GET requests
			AllowedOrigin: []string{"*"},   // Allow any origin (use specific one in prod)
			ExposeHeader:  []string{"ETag", "Content-Type"},
		},
	}))
	if test != nil {
		fmt.Println("Failed to set bucket CORS:", test)
	}

	return &FileService{
		minioClient: minioClient,
		bucketName:  bucketName,
		cache:       fileUrlCache,
	}, nil
}

func (f *FileService) UploadFile(objectName string, file io.Reader, fileSize int64, mimeType string) (string, error) {
	_, err := f.minioClient.PutObject(
		context.Background(),
		f.bucketName,
		objectName,
		file,
		fileSize,
		minio.PutObjectOptions{ContentType: mimeType, CacheControl: "public, max-age=31536000"},
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", f.minioClient.EndpointURL(), f.bucketName, objectName), nil
}

func (f *FileService) GetFileUrl(objectName string) (string, error) {
	cachedUrl, err := f.cache.Get(objectName)
	if err == nil {
		return cachedUrl, nil
	}
	url, err := f.minioClient.PresignedGetObject(context.Background(), f.bucketName, objectName, time.Duration(36)*time.Hour, nil)
	if err != nil {
		return "", err
	}
	f.cache.Set(objectName, url.String(), time.Duration(36)*time.Hour)
	return url.String(), nil
}

package s3

import (
	"context"
	"fmt"
	"log"
	"teach-stack/types"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	userPhotosBucket = "user-photos"
)

type MinioProvider struct {
	minioAuthData
	client *minio.Client
}

type minioAuthData struct {
	url      string
	user     string
	password string
	token    string
	ssl      bool
}

func NewMinioProvider(minioURL string, minioUser string, minioPassword string, ssl bool) *MinioProvider {
	return &MinioProvider{
		minioAuthData: minioAuthData{
			password: minioPassword,
			url:      minioURL,
			user:     minioUser,
			ssl:      ssl,
		}}
}

func (m *MinioProvider) Connect() error {
	var err error
	m.client, err = minio.New(m.url, &minio.Options{
		Creds:  credentials.NewStaticV4(m.user, m.password, ""),
		Secure: m.ssl,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func (m *MinioProvider) Upload(ctx context.Context, img types.Image) (string, error) {
	if m.client == nil {
		return "", fmt.Errorf("MinIO client не инициализирован")
	}

	imageName := "image/" + uuid.New().String()

	_, err := m.client.PutObject(
		ctx,
		userPhotosBucket,
		imageName,
		img.Payload,
		img.Size,
		minio.PutObjectOptions{
			ContentType: "image/png",
		},
	)
	if err != nil {
		return "", err
	}
	return imageName, nil
}

func (m *MinioProvider) GetTemplUrl(ctx context.Context, fileName string) (string, error) {
	// Ensure the bucket exists
	exists, err := m.client.BucketExists(ctx, userPhotosBucket)
	if err != nil {
		return "", fmt.Errorf("error checking bucket existence: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("bucket %s does not exist", userPhotosBucket)
	}

	// Check if the object exists
	_, err = m.client.StatObject(ctx, userPhotosBucket, fileName, minio.StatObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("error checking object existence: %w", err)
	}

	// Generate presigned URL with longer expiration
	url, err := m.client.PresignedGetObject(
		ctx,
		userPhotosBucket,
		fileName,
		time.Hour*24, // 24 hours
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("error generating presigned URL: %w", err)
	}

	// Log the generated URL for debugging
	log.Printf("Generated presigned URL for %s: %s", fileName, url.String())

	return url.String(), nil
}

func (m *MinioProvider) DownloadImage(ctx context.Context, bucketName string, fileName string) (*minio.Object, error) {
	obj, err := m.client.GetObject(
		ctx,
		bucketName,
		fileName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

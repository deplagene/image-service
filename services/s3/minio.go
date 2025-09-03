package s3

import (
	"context"
	"deplagene/image-service/types"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Provider struct {
	authData
	client *minio.Client
}

type authData struct {
	url      string
	user     string
	password string
	token    string
	ssl      bool
}

func NewProvider(user, password, token string, url string, ssl bool) *Provider {
	return &Provider{
		authData: authData{
			url:      url,
			user:     user,
			password: password,
			token:    token,
			ssl:      ssl,
		},
	}
}

func (p *Provider) Connect() error {
	const op = "s3.Provider.Connect"

	var err error
	p.client, err = minio.New(p.url, &minio.Options{
		Creds:  credentials.NewStaticV4(p.user, p.password, ""),
		Secure: p.ssl,
	})
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (p *Provider) Upload(ctx context.Context, img types.Image) (string, error) {
	const op = "s3.Provider.Upload"

	if p.client == nil {
		return "", fmt.Errorf("%s: client is not connected", op)
	}

	imageName := generateUploadImageName()

	_, err := p.client.PutObject(
		ctx,
		types.UserImagesBucketName,
		imageName,
		img.Payload,
		img.Size,
		minio.PutObjectOptions{
			ContentType: "image/png",
		},
	)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return imageName, nil
}

func generateUploadImageName() string {
	const serviceName = "image-service"
	uuidStr := strings.ReplaceAll(uuid.NewString(), "-", "")

	return fmt.Sprintf("%s-%s", serviceName, uuidStr)
}

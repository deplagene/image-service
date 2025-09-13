package types

import (
	"context"
	"io"

	"github.com/google/uuid"
)

const (
	MinConfidenceThreshold = 85.0
	ImgPath                = "assets/nsfw-content.jpg"
	UserImagesBucketName   = "user-images"
)

type NsfwResult struct {
	FileName             string  `json:"file_name"`
	IsNsfw               bool    `json:"is_nsfw"`
	ConfidencePercentage float64 `json:"confidence_percentage"`
}

type Image struct {
	Id          uuid.UUID
	PayloadName string
	Payload     io.Reader
	Size        int64
	Url         string
}

type ImageService interface {
	CheckNsfw(ctx context.Context, image Image) (NsfwResult, error)
}

type ImageStore interface {
	Create(ctx context.Context, url string) error
	GetByUrl(ctx context.Context, url string) (Image, error)
}

type S3Storage interface {
	Connect() error
	Upload(ctx context.Context, image Image) (string, error)
}

type Broker interface {
	Publish(ctx context.Context, topic string, data string) error
	Consume(ctx context.Context, topic string, groupId string) error
}

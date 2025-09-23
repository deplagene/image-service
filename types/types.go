package types

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const (
	MinConfidenceThreshold = 85.0
	ImgPath                = "assets/nsfw-content.jpg"
	UserImagesBucketName   = "user-images"
	MaxUploadSize          = 10 << 20 // 10 Mb
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
	Upload(ctx context.Context) error
	CheckImageForNsfwByUrl(ctx context.Context, url string) (NsfwResult, error)
	CheckImageForNsfwByUpload(ctx context.Context, image Image) (NsfwResult, error)
}

type ImageStore interface {
	Create(ctx context.Context, url string) error
	GetByUrl(ctx context.Context, url string) (Image, error)
}

type S3Storage interface {
	Connect() error
	Upload(ctx context.Context, image Image) (string, error)
}

type Producer interface {
	Publish(ctx context.Context, msg kafka.Message) error
	Close() error
}

type Consumer interface {
	Consume(handler func(ctx context.Context, msg kafka.Message) error)
	Close() error 
}

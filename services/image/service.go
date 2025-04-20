package image

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"teach-stack/types"

	"github.com/disintegration/imaging"
)

const (
	SharpenFilter       = "sharpen"
	BlackAndWhiteFilter = "bw"
)

type Service struct {
	store types.ImageStore
	s3    types.S3_Storage
}

func NewService(store types.ImageStore, s3 types.S3_Storage) *Service {
	return &Service{
		store: store,
		s3:    s3,
	}
}

func (s *Service) AddFilters(task types.PhotoProcessingTask) error {
	obj, err := s.s3.DownloadImage(
		context.Background(),
		"user-photos",
		task.FileName,
	)
	if err != nil {
		slog.Error("Cannot download image from S3.")
		return err
	}

	defer obj.Close()

	info, err := obj.Stat()
	if err != nil {
		slog.Error("Cannot get object info", "error", err)
		return err
	}
	slog.Info("Object info", "size", info.Size, "contentType", info.ContentType)

	img, err := imaging.Decode(obj)
	if err != nil {
		slog.Error("Cannot decode object to image.")
		return err
	}

	processedImage := img
	for _, f := range task.Filters {
		switch f {
		case BlackAndWhiteFilter:
			processedImage = imaging.Grayscale(processedImage)
		case SharpenFilter:
			processedImage = imaging.Sharpen(processedImage, 1.5)
		default:
			slog.Warn("Unknown filter", "filter", f)
		}
	}

	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, processedImage, imaging.PNG); err != nil {
		slog.Error("Cannot encode processed image", "error", err, "fileName", task.FileName)
		return fmt.Errorf("failed to encode processed image: %w", err)
	}

	newImage := types.Image{
		Payload: buf,
		Size:    int64(buf.Len()),
	}

	uploadedName, err := s.s3.Upload(context.Background(), newImage)
	if err != nil {
		return err
	}

	// Get presigned URL for the uploaded image
	url, err := s.s3.GetTemplUrl(context.Background(), uploadedName)
	if err != nil {
		slog.Error("Failed to get presigned URL", "error", err)
		return err
	}

	if err := s.store.Update(task.Id, url); err != nil {
		slog.Error("Cannot create image", "error", err, "fileName", task.FileName)
		return fmt.Errorf("failed to create image: %w", err)
	}

	return nil
}

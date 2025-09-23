package images

import (
	"context"
	"deplagene/image-service/types"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func Upload(ctx context.Context) error {
	const op = "image.service.Upload"

	return nil
}

// todo: реализовать логику

func (s *Service) CheckImageForNsfwByUpload(ctx context.Context, image types.Image) (types.NsfwResult, error) {
	const op = "image.service.CheckNsfw"

	return types.NsfwResult{}, nil

}

// todo : реализовать логику

func (s *Service) CheckImageForNsfwByUrl(ctx context.Context, url string) (types.NsfwResult, error) {
	const op = "image.service.CheckNsfwByUrl"

	return types.NsfwResult{}, nil
}

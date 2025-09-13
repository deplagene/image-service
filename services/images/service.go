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

// todo: реализовать логику

func (s *Service) CheckNsfw(ctx context.Context, image types.Image) (types.NsfwResult, error) {
	const op = "image.service.CheckNsfw"

	return types.NsfwResult{}, nil
}

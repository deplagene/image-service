package images

import (
	"context"
	"deplagene/image-service/types"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	imageDatabase   = "image-service-db"
	imageCollection = "image"
)

type Store struct {
	cl *mongo.Client
}

func NewStore(cl *mongo.Client) *Store {
	return &Store{
		cl: cl,
	}
}

func (s *Store) Create(ctx context.Context, url string) error {
	const op = "images.Store.Create"

	db := s.cl.Database(imageDatabase)
	imgCollection := db.Collection(imageCollection)

	img := types.Image{
		Url: url,
	}

	res, err := imgCollection.InsertOne(ctx, img)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Printf("id: %v", res.InsertedID)

	return nil
}

func (s *Store) GetByUrl(ctx context.Context, url string) (types.Image, error) {
	const op = "images.Store.GetByUrl"

	db := s.cl.Database(imageDatabase)
	imgCollection := db.Collection(imageCollection)

	var img types.Image
	err := imgCollection.FindOne(ctx, types.Image{Url: url}).Decode(&img)
	if err != nil {
		return types.Image{}, fmt.Errorf("%s: %w", op, err)
	}

	return img, nil
}

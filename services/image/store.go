package image

import (
	"database/sql"
	"teach-stack/types"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Create(i types.Image) error {
	_, err := s.db.Exec(createImageQuery, i.Id, i.PayloadName, i.Size, i.Url)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetByUrl(url string) (*types.Image, error) {
	rows, err := s.db.Query(getImageByUrlQuery, url)
	if err != nil {
		return nil, err
	}
	i := new(types.Image)
	for rows.Next() {
		i, err = scanRowIntoImages(rows)
		if err != nil {
			return nil, err
		}
	}

	return i, nil
}

func(s *Store) GetById(id uuid.UUID) (*types.Image, error) {
	rows, err := s.db.Query(getImageByIdQuery, id)
	if err != nil {
		return nil, err
	}

	i := new(types.Image)
	for rows.Next() {
		i, err = scanRowIntoImages(rows)
		if err != nil {
			return nil, err
		}
	}

	return i, err
}

func (s *Store) Update(id uuid.UUID, url string) error {
	rows, err := s.db.Query(updateImageQuery, url, id)
	if err != nil {
		return err
	}
	i := new(types.Image)
	for rows.Next() {
		i, err = scanRowIntoImages(rows)
		if err != nil {
			return err
		}
	}

	_ = i

	return nil
}

func scanRowIntoImages(rows *sql.Rows) (*types.Image, error) {
	photo := new(types.Image)

	err := rows.Scan(
		&photo.Id,
		&photo.PayloadName,
		&photo.Size,
		&photo.Url,
	)
	if err != nil {
		return nil, err
	}

	return photo, nil
}

package types

import (
	"encoding/json"
	"fmt"
)

// todo: add fields
type ImageUploadEvent struct {
}

func (i *ImageUploadEvent) Marshal() ([]byte, error) {
	const op = "types.ImageUploadEvent.Marshal"

	data, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}

func UnmarshalUploadedImage(data []byte) (*ImageUploadEvent, error) {
	const op = "types.ImageUploadEvent.Unmarshal"

	var event ImageUploadEvent

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &event, nil
}

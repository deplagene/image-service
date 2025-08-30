package types

const (
	MinConfidenceThreshold = 85.0
	ImgPath                = "assets/nsfw-content-2.jpg"
)

type NsfwResult struct {
	FileName             string  `json:"file_name"`
	IsNsfw               bool    `json:"is_nsfw"`
	ConfidencePercentage float64 `json:"confidence_percentage"`
}

package images

import (
	"deplagene/image-service/types"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r *http.ServeMux) {
}

func (h *Handler) uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxUploadSize)

	if err := r.ParseMultipartForm(types.MaxUploadSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// todo: use const
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = fileHeader

	defer file.Close()
}

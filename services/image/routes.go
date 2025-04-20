package image

import (
	_ "context"
	"encoding/json"
	"log/slog"
	"net/http"
	"teach-stack/types"
	"teach-stack/utils"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	store   types.ImageStore
	s3      types.S3_Storage
	broker  types.MessageBroker
	service types.ImageService
}

func NewHandler(store types.ImageStore, s3 types.S3_Storage, broker types.MessageBroker, service types.ImageService) *Handler {
	return &Handler{
		store:   store,
		s3:      s3,
		broker:  broker,
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /upload", h.uploadHandler)
	r.HandleFunc("GET /images/{id}", h.getResultHandler)
}

func (h *Handler) uploadHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	rec := &utils.StatusRecorder{ResponseWriter: w, Status: http.StatusOK}
	defer func() {
		observeRequest(time.Since(start), rec.Status)
	}()
	
	const maxUploadSize = 10 << 20 // 10 Mb
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file form:"+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	object := types.Image{
		PayloadName: fileHeader.Filename,
		Payload:     file,
		Size:        fileHeader.Size,
	}

	fileName, err := h.s3.Upload(r.Context(), object)
	if err != nil {
		http.Error(w, "Error uploading an image to storage"+err.Error(), http.StatusInternalServerError)
		return
	}

	img := types.Image{
		Id:          uuid.New(),
		PayloadName: fileName,
		Size:        fileHeader.Size,
		Url:         "",
	}

	if err := h.store.Create(img); err != nil {
		http.Error(w, "Failed to create photo in store: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.broker.CreateExchange("image", amqp.ExchangeFanout, true); err != nil {
		http.Error(w, "Error create exchange"+err.Error(), http.StatusBadRequest)
		return
	}

	task := types.PhotoProcessingTask{
		Id:       img.Id,
		FileName: fileName,
		Filters:  []string{"bw"},
	}

	body, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Failed to marshal task"+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.broker.Send(r.Context(), "image", "", amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		http.Error(w, "Failed to send task to queue"+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(map[string]string{
		"photo_id": img.Id.String(),
	})
	if err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getResultHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Error parsing id from query", http.StatusBadRequest)
		return
	}

	image, err := h.store.GetById(uuid.MustParse(id))
	if err != nil {
		slog.Error("Error getting image from store", "error", err)
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Retrieved image from store", "image", image)

	// url, err := h.s3.GetTemplUrl(context.Background(), image.PayloadName)
	// if err != nil {
	// 	slog.Error("Error getting temporary URL", "error", err, "fileName", image.PayloadName)
	// 	http.Error(w, "Error get url: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// slog.Info("Generated temporary URL", "url", url)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{
		"url": image.Url,
	})
	if err != nil {
		slog.Error("Error encoding response", "error", err)
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

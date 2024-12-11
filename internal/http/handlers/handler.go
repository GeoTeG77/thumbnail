package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"thumbnail/internal/models"
	"thumbnail/internal/service"
)

type ThumbnailHandler struct {
	Service *service.ThumbnailService
}

func NewThumbnailHandler(ser *service.ThumbnailService) *ThumbnailHandler {
	return &ThumbnailHandler{Service: ser}
}

func (h *ThumbnailHandler) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	request := make([]*models.ThumbnailRequest, 0)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("Failed to decode JSON request", "error", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	slog.Info("Received request", "count", len(request))
	ctx := r.Context()
	if len(request) == 1 {
		thumbnail, err := h.Service.GetThumbnail(ctx, request[0].URL)
		if err != nil {
			slog.Error("Cache Miss", "url", request[0].URL)
			parts := strings.Split(request[0].URL, "v=")
			if len(parts) < 2 {
				http.Error(w, "Invalid YouTube URL format", http.StatusBadRequest)
				return
			}

			thumbnailURL := "https://img.youtube.com/vi/" + parts[1] + "/0.jpg"
			resp, err := http.Get(thumbnailURL)
			if err != nil {
				slog.Error("Failed to fetch thumbnail from YouTube", "url", thumbnailURL, "error", err)
				http.Error(w, "Failed to fetch thumbnail", http.StatusBadRequest)
				return
			}

			defer resp.Body.Close()

			imgBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("Failed to read thumbnail response body", "error", err)
				http.Error(w, "Failed to read thumbnail", http.StatusInternalServerError)
				return
			}

			thumbnail = &models.ThumbnailResponse{
				URL:     request[0].URL,
				Preview: imgBytes,
			}
		}
		slog.Info("Successfully generated thumbnail", "url", request[0].URL)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(thumbnail); err != nil {
			slog.Error("Failed to encode JSON response", "error", err)
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		}
		return
	}

	response := make([]*models.ThumbnailResponse, 0, len(request))
	for idx, req := range request {
		thumbnail, err := h.Service.GetThumbnail(ctx, req.URL)
		if err != nil {
			slog.Error("Failed to generate thumbnail", "url", req.URL, "error", err)
			continue
		}
		slog.Info("Successfully generated thumbnail", "url", req.URL)
		response[idx] = thumbnail
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}
	slog.Info("Successfully processed request", "count", len(request))
	return
}

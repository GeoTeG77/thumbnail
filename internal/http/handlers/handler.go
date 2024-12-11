package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"thumbnail/internal/models"
	"thumbnail/internal/service"
)

type ThumbnailHandler struct {
	Service *service.ThumbnailService
}

func NewThumbnailHandler(service *service.ThumbnailService) *ThumbnailHandler {
	return &ThumbnailHandler{Service: service}
}

func (h *ThumbnailHandler) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	var request []*models.ThumbnailRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("Failed to decode JSON request", "error", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	slog.Info("Received request", "count", len(request))
	ctx := r.Context()

	if len(request) == 1 {
		h.handleSingleThumbnailRequest(w, r, request[0])
		return
	}

	h.handleMultipleThumbnailRequests(w, ctx, request)
}

func (h *ThumbnailHandler) handleSingleThumbnailRequest(w http.ResponseWriter, r *http.Request, req *models.ThumbnailRequest) {
	thumbnail, err := h.Service.GetThumbnail(r.Context(), req.URL)
	if err == nil {
		slog.Info("Thumbnail Taken from Cache")
	}
	if err != nil {
		slog.Error("Cache miss for URL", "url", req.URL)
		thumbnail = h.fetchThumbnailFromYouTube(w, req)
		if thumbnail == nil {
			return
		}
		err := h.Service.SetThumbnail(r.Context(), thumbnail.URL, thumbnail.Preview)
		if err != nil {
			slog.Error("Failed to save Thumbnail in Cache")
		}
		if err == nil {
			slog.Info("Saved thumbnail in Cache for url:", "url", req.URL)
		}
	}
	slog.Info("Successfully take thumbnail from YouTube", "url", req.URL)
	h.writeImageResponse(w, thumbnail.Preview)
}

func (h *ThumbnailHandler) handleMultipleThumbnailRequests(w http.ResponseWriter, ctx context.Context, request []*models.ThumbnailRequest) {
	response := make([]*models.ThumbnailResponse, len(request))
	for idx, req := range request {
		thumbnail, err := h.Service.GetThumbnail(ctx, req.URL)
		if err == nil {
			slog.Info("Thumbnail Taken from Cache")
		}
		if err != nil {
			slog.Error("Cache miss for URL", "url", req.URL)
			thumbnail = h.fetchThumbnailFromYouTube(w, req)
			if thumbnail == nil {
				return
			}
			err := h.Service.SetThumbnail(ctx, thumbnail.URL, thumbnail.Preview)
			if err != nil {
				slog.Error("Failed to save Thumbnail in Cache")
			}
			if err == nil {
				slog.Info("Saved thumbnail in Cache for url:", "url", req.URL)
			}
		}
		slog.Info("Successfully take thumbnail from YouTube", "url", req.URL)
		response[idx] = thumbnail
	}

	h.writeImagesResponse(w, response)
}

func (h *ThumbnailHandler) fetchThumbnailFromYouTube(w http.ResponseWriter, req *models.ThumbnailRequest) *models.ThumbnailResponse {
	parts := strings.Split(req.URL, "v=")
	if len(parts) < 2 {
		http.Error(w, "Invalid YouTube URL format", http.StatusBadRequest)
		return nil
	}

	thumbnailURL := os.Getenv("THUMBNAIL_URL_P1") + parts[1] + os.Getenv("THUMBNAIL_URL_P2")

	resp, err := http.Get(thumbnailURL)
	if err != nil {
		slog.Error("Failed to fetch thumbnail from YouTube", "url", thumbnailURL, "error", err)
		http.Error(w, "Failed to fetch thumbnail", http.StatusBadRequest)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Received non-OK status code", "url", thumbnailURL, "status", resp.StatusCode)
		http.Error(w, "Failed to fetch thumbnail", http.StatusInternalServerError)
		return nil
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image/jpeg") {
		slog.Error("Expected JPEG image, but got", "content-type", resp.Header.Get("Content-Type"))
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return nil
	}

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read thumbnail response body", "error", err)
		http.Error(w, "Failed to read thumbnail", http.StatusInternalServerError)
		return nil
	}

	return &models.ThumbnailResponse{
		URL:     req.URL,
		Preview: imgBytes,
	}
}

func (h *ThumbnailHandler) writeImageResponse(w http.ResponseWriter, imgBytes []byte) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(imgBytes)))
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(imgBytes); err != nil {
		slog.Error("Failed to write image response", "error", err)
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
	}
}

func (h *ThumbnailHandler) writeImagesResponse(w http.ResponseWriter, response []*models.ThumbnailResponse) {
	w.Header().Set("Content-Type", "multipart/related; boundary=myBoundary")
	w.WriteHeader(http.StatusOK)

	for idx, thumbnail := range response {
		partHeader := fmt.Sprintf("--myBoundary\r\nContent-Type: image/jpeg\r\nContent-Disposition: inline; filename=\"thumbnail_%d.jpg\"\r\nContent-Length: %d\r\n\r\n", idx, len(thumbnail.Preview))

		if _, err := w.Write([]byte(partHeader)); err != nil {
			slog.Error("Failed to write part header", "error", err)
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(thumbnail.Preview); err != nil {
			slog.Error("Failed to write image data", "error", err)
			http.Error(w, "Failed to generate response", http.StatusInternalServerError)
			return
		}
	}

	if _, err := w.Write([]byte("\r\n--myBoundary--\r\n")); err != nil {
		slog.Error("Failed to write boundary end", "error", err)
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
	}
}
